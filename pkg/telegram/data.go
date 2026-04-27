package telegram

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/hyperized/dsmr/pkg/obis"
)

// Data holds the parsed result of a single OBIS data line.
type Data struct {
	raw         string
	value       []string
	typedValues []any
	reference   *obis.Reference
}

// obisLineRE captures a canonical DSMR OBIS identifier and the parenthesised
// value list. The B/C separator is `.` per IEC 62056-61; lines using `:`
// (seen in pre-DSMR-5 telegrams) are deliberately rejected so they fall
// through to the parser's noise path instead of producing unknown lookups.
var obisLineRE = regexp.MustCompile(`([0-9]-[0-9]:[0-9]+\.[0-9]+\.[0-9]+)\((.*)`)

// NewData parses an OBIS data line. It returns an error if the line does not match
// the expected OBIS pattern or if the identifier is not in the reference registry.
func NewData(line string) (*Data, error) {
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

	d.typedValues = parseTyped(value, reference)

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
func parseTyped(values []string, ref *obis.Reference) []any {
	f := ref.Format()
	fs := f.FormatString()

	// GenericProfile (e.g. power-failure event log) has a complex layout that
	// we don't decode into typed values; the raw values list is preserved so
	// callers can interpret it if they need to.
	if f.Class() == cosem.GenericProfile {
		return nil
	}

	if len(values) == 0 || fs == "" {
		return nil
	}

	var typed []any

	switch {
	case strings.HasPrefix(fs, "F"):
		// FloatingPoint: use the format string (e.g. "F9(3,3)") directly.
		if v, err := cosem.NewFloatingPoint(values[0], cosem.WithFormat(fs)); err == nil {
			typed = append(typed, v)
		} else {
			slog.Debug("FloatingPoint parse error", "value", values[0], "format", fs, "err", err)
		}

	case fs == "TST":
		if f.Class() == cosem.ExtendedRegister {
			// M-Bus extended register: values[0]=capture timestamp, values[1]=metered value.
			if len(values) >= 1 {
				if ts, err := cosem.NewTimestamp(values[0]); err == nil {
					typed = append(typed, ts)
				} else {
					slog.Debug("M-Bus timestamp parse error", "value", values[0], "err", err)
				}
			}
			if len(values) >= 2 {
				// Metered value kept raw; the parser converts it to float when
				// emitting the M-Bus gauge.
				typed = append(typed, values[1])
			}
		} else {
			if ts, err := cosem.NewTimestamp(values[0]); err == nil {
				typed = append(typed, ts)
			} else {
				slog.Debug("timestamp parse error", "value", values[0], "err", err)
			}
		}

	case strings.HasPrefix(fs, "I"):
		// Integer with fixed digit count, e.g. "I3".
		length, _ := strconv.Atoi(fs[1:])
		if v, err := cosem.NewInteger(values[0], cosem.WithIntegerLength(length)); err == nil {
			typed = append(typed, v)
		} else {
			slog.Debug("integer parse error", "value", values[0], "format", fs, "err", err)
		}

	default:
		// S-format strings and any unknown format: store the raw string value.
		typed = append(typed, values[0])
	}

	return typed
}
