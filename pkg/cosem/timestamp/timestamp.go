// Package timestamp implements the COSEM TST (timestamp) type.
// DSMR timestamps use the format YYMMDDhhmmssX where X is 'S' (summer/DST)
// or 'W' (winter), interpreted in the Europe/Amsterdam time zone.
package timestamp

import (
	"fmt"
	"time"

	"github.com/hyperized/dsmr/pkg/cosem"
)

// OptionFunc configures a Timestamp.
type OptionFunc func(t *Timestamp)

// DSTFlag represents the daylight saving time indicator in a DSMR timestamp.
type DSTFlag byte

// DST flag values as they appear in DSMR timestamps (13th character).
const (
	Winter DSTFlag = 'W'
	Summer DSTFlag = 'S'
)

// Timestamp holds a parsed DSMR timestamp together with its DST flag.
type Timestamp struct {
	tag   cosem.Cosem
	dst   DSTFlag
	value time.Time
}

// New parses a DSMR timestamp string of the form YYMMDDhhmmssX where X is 'S'
// (summer/DST) or 'W' (winter). The time is interpreted in the Europe/Amsterdam
// time zone, falling back to UTC if the zone data is unavailable.
func New(value string, options ...OptionFunc) (*Timestamp, error) {
	t := &Timestamp{
		tag: cosem.OctetString,
	}

	for _, o := range options {
		o(t)
	}

	if len(value) != 13 {
		return nil, fmt.Errorf("timestamp must be 13 characters, found %d", len(value))
	}

	dstChar := value[12]
	if dstChar != byte(Winter) && dstChar != byte(Summer) {
		return nil, fmt.Errorf("invalid DST flag %q, expected 'S' or 'W'", dstChar)
	}

	t.dst = DSTFlag(dstChar)

	loc, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		loc = time.UTC
	}

	parsed, err := time.ParseInLocation(cosem.DateTimeFormat, value[:12], loc)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp %q: %w", value[:12], err)
	}

	t.value = parsed
	return t, nil
}

// WithTag overrides the default COSEM tag for this timestamp.
func WithTag(tag cosem.Cosem) OptionFunc {
	return func(t *Timestamp) {
		t.tag = tag
	}
}

// Value returns the parsed time in the Europe/Amsterdam time zone.
func (t *Timestamp) Value() time.Time {
	return t.value
}

// DST returns the raw DST flag byte ('W' or 'S') from the timestamp string.
func (t *Timestamp) DST() DSTFlag {
	return t.dst
}

// IsSummer reports whether the timestamp uses summer time (DST active).
func (t *Timestamp) IsSummer() bool {
	return t.dst == Summer
}
