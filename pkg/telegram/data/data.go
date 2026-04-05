// Package data parses a single DSMR OBIS data line (e.g. "1-0:1.8.1(123456.789*kWh)")
// into a structured Data value.  It looks up the OBIS reference, splits
// multi-value fields, strips units, and converts raw strings into their
// concrete COSEM types (FloatingPoint, Integer, Timestamp, …).
package data

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hyperized/dsmr/pkg/cosem"
	fp "github.com/hyperized/dsmr/pkg/cosem/floating_point"
	"github.com/hyperized/dsmr/pkg/cosem/integer"
	"github.com/hyperized/dsmr/pkg/cosem/timestamp"
	"github.com/hyperized/dsmr/pkg/obis"
)

// PowerFailureEvent holds a single entry from the power failure event log
// (OBIS 1-0:99.97.0).
type PowerFailureEvent struct {
	Timestamp time.Time
	Duration  time.Duration
}

// Data holds the parsed result of a single OBIS data line.
type Data struct {
	raw         string
	value       []string
	typedValues []any
	events      []PowerFailureEvent
	reference   *obis.Reference
}

var obisLineRE = regexp.MustCompile(`([0-9]-[0-9]:[0-9]+[\.:][0-9]+\.[0-9]+)\((.*)`)

// New parses an OBIS data line. It returns an error if the line does not match
// the expected OBIS pattern or if the identifier is not in the reference registry.
func New(line string) (*Data, error) {
	match := obisLineRE.FindStringSubmatch(strings.TrimSpace(line))
	if len(match) != 3 {
		return nil, errors.New("data: could not parse line")
	}

	// The captured group includes everything from the first '(' to end of line.
	// Trim the trailing ')' so splitValues receives clean inner content.
	raw := strings.TrimSuffix(match[2], ")")

	reference := obis.New(match[1])
	if reference == nil {
		return nil, fmt.Errorf("data: unknown OBIS reference %q", match[1])
	}

	value := splitValues(raw)
	if len(value) == 0 || value[0] == "" {
		return nil, errors.New("data: no value could be parsed from this line")
	}

	d := &Data{
		raw:       raw,
		value:     value,
		reference: reference,
	}

	d.typedValues, d.events = parseTyped(value, reference)

	return d, nil
}

// Identifier returns the OBIS identifier string for this data line (e.g. "1-0:1.8.1").
func (d *Data) Identifier() string {
	return d.reference.Identifier()
}

// MetricName returns the Prometheus metric name for this data line, or "" if none.
func (d *Data) MetricName() string {
	return d.reference.MetricName()
}

// Values returns the raw string values extracted from the parenthesised fields.
func (d *Data) Values() []string {
	return d.value
}

// TypedValues returns the parsed COSEM-typed values for this data line.
func (d *Data) TypedValues() []any {
	return d.typedValues
}

// Events returns the parsed power-failure event log entries (OBIS 1-0:99.97.0).
func (d *Data) Events() []PowerFailureEvent {
	return d.events
}

func (d *Data) String() string {
	return fmt.Sprintf("%s: %s%s", d.reference.Name(), d.value[0], d.reference.Unit())
}

// splitValues splits a multi-value DSMR string on ")(" boundaries and strips
// the unit suffix (everything from '*' onward) from each element.
func splitValues(raw string) []string {
	parts := strings.Split(raw, ")(")
	r := make([]string, len(parts))
	for i, v := range parts {
		r[i] = strings.SplitN(v, "*", 2)[0]
	}
	return r
}

// parseTyped attempts to parse raw string values into their COSEM types using
// the format specification from the OBIS reference.
func parseTyped(values []string, ref *obis.Reference) ([]any, []PowerFailureEvent) {
	f := ref.Format()
	fs := f.FormatString()

	// Power failure event log uses GenericProfile class with a complex layout.
	if f.Class() == cosem.GenericProfile {
		return nil, parsePowerFailureLog(values)
	}

	if len(values) == 0 || fs == "" {
		return nil, nil
	}

	var typed []any

	switch {
	case strings.HasPrefix(fs, "F"):
		// FloatingPoint: use the format string (e.g. "F9(3,3)") directly.
		if v, err := fp.New(values[0], fp.WithFormat(fs)); err == nil {
			typed = append(typed, v)
		} else {
			slog.Warn("FloatingPoint parse error", "value", values[0], "format", fs, "err", err)
		}

	case fs == "TST":
		if f.Class() == cosem.ExtendedRegister {
			// M-Bus extended register: values[0]=capture timestamp, values[1]=metered value.
			if len(values) >= 1 {
				if ts, err := timestamp.New(values[0]); err == nil {
					typed = append(typed, ts)
				} else {
					slog.Warn("M-Bus timestamp parse error", "value", values[0], "err", err)
				}
			}
			if len(values) >= 2 {
				// Metered value kept raw until device type is resolved (Phase 4.3).
				typed = append(typed, values[1])
			}
		} else {
			if ts, err := timestamp.New(values[0]); err == nil {
				typed = append(typed, ts)
			} else {
				slog.Warn("timestamp parse error", "value", values[0], "err", err)
			}
		}

	case strings.HasPrefix(fs, "I"):
		// Integer with fixed digit count, e.g. "I3".
		length, _ := strconv.Atoi(fs[1:])
		if v, err := integer.New(values[0], integer.WithLength(length)); err == nil {
			typed = append(typed, v)
		} else {
			slog.Warn("integer parse error", "value", values[0], "format", fs, "err", err)
		}

	default:
		// S-format strings and any unknown format: store the raw string value.
		typed = append(typed, values[0])
	}

	return typed, nil
}

// parsePowerFailureLog parses the value list of OBIS 1-0:99.97.0.
//
// After splitValues the layout is:
//
//	[count, "0-0:96.7.19", TST1, duration1_s, TST2, duration2_s, …]
func parsePowerFailureLog(values []string) []PowerFailureEvent {
	// Minimum: count + obis-ref + one (TST, duration) pair.
	if len(values) < 4 {
		return nil
	}

	var events []PowerFailureEvent

	// values[0] = event count (informational); values[1] = OBIS ref (skip).
	for i := 2; i+1 < len(values); i += 2 {
		ts, err := timestamp.New(values[i])
		if err != nil {
			slog.Warn("power failure timestamp parse error", "value", values[i], "err", err)
			continue
		}
		secs, err := strconv.ParseInt(values[i+1], 10, 64)
		if err != nil {
			slog.Warn("power failure duration parse error", "value", values[i+1], "err", err)
			continue
		}
		events = append(events, PowerFailureEvent{
			Timestamp: ts.Value(),
			Duration:  time.Duration(secs) * time.Second,
		})
	}

	return events
}
