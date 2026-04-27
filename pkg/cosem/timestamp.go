package cosem

import (
	"fmt"
	"time"
)

// DSTFlag represents the daylight saving time indicator in a DSMR timestamp.
type DSTFlag byte

// DST flag values as they appear in DSMR timestamps (13th character).
const (
	Winter DSTFlag = 'W'
	Summer DSTFlag = 'S'
)

// amsterdam is the cached Europe/Amsterdam timezone used to interpret DSMR
// timestamps. Falls back to UTC if the zoneinfo database is unavailable.
// Cached at init because time.LoadLocation is not internally memoised and
// timestamp parsing is on the per-line hot path.
var amsterdam = loadLocationOrUTC("Europe/Amsterdam")

func loadLocationOrUTC(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return loc
}

// Timestamp holds a parsed DSMR timestamp together with its DST flag.
type Timestamp struct {
	dst   DSTFlag
	value time.Time
}

// NewTimestamp parses a DSMR timestamp string of the form YYMMDDhhmmssX where
// X is 'S' (summer/DST) or 'W' (winter). The time is interpreted in the
// Europe/Amsterdam time zone, falling back to UTC if the zone data is
// unavailable.
func NewTimestamp(value string) (*Timestamp, error) {
	if len(value) != 13 {
		return nil, fmt.Errorf("timestamp must be 13 characters, found %d", len(value))
	}

	dstChar := value[12]
	if dstChar != byte(Winter) && dstChar != byte(Summer) {
		return nil, fmt.Errorf("invalid DST flag %q, expected 'S' or 'W'", dstChar)
	}

	parsed, err := time.ParseInLocation(DateTimeFormat, value[:12], amsterdam)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp %q: %w", value[:12], err)
	}

	return &Timestamp{dst: DSTFlag(dstChar), value: parsed}, nil
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
