// Package floating_point implements the COSEM FloatingPoint type.
// Values are decimal strings that are validated against optional min/max
// decimal-digit constraints before being stored as float64.
package floating_point

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/shopspring/decimal"
)

// OptionFunc configures a FloatingPoint value.
type OptionFunc func(f *FloatingPoint)

// FloatingPoint holds a parsed COSEM floating-point value with its type metadata.
type FloatingPoint struct {
	tag             cosem.Cosem
	length          int32
	minimumDecimals int
	maximumDecimals int
	value           float64
}

// New parses value as a decimal floating-point number. The value must contain
// only digits and at most one dot. Decimal-digit count must satisfy any
// min/max constraints set via options.
func New(value string, options ...OptionFunc) (*FloatingPoint, error) {
	f := &FloatingPoint{
		tag:             cosem.Float32,
		length:          0,
		minimumDecimals: 0,
		maximumDecimals: 0,
		value:           0,
	}

	for _, o := range options {
		o(f)
	}

	if !isValidString(value) {
		return nil, errors.New("the provided string can only contain digits and dots")
	}
	digits, err := digitsAfterDot(value)
	if err != nil {
		return nil, err
	}

	if digits < f.minimumDecimals || digits > f.maximumDecimals {
		return nil, fmt.Errorf("decimal point incorrect, expected to be between %d and %d, found %d",
			f.minimumDecimals, f.maximumDecimals, digits,
		)
	}

	f.value = stringToFloat64(value, f.length)

	return f, nil
}

// WithTag overrides the default COSEM tag for this floating-point value.
func WithTag(tag cosem.Cosem) OptionFunc {
	return func(f *FloatingPoint) {
		f.tag = tag
	}
}

// WithMinimumDecimals sets the minimum number of digits after the decimal point.
func WithMinimumDecimals(minimum int) OptionFunc {
	return func(f *FloatingPoint) {
		f.minimumDecimals = minimum
	}
}

// WithMaximumDecimals sets the maximum number of digits after the decimal point.
func WithMaximumDecimals(maximum int) OptionFunc {
	return func(f *FloatingPoint) {
		f.maximumDecimals = maximum
	}
}

// WithLength sets the total field width (used for decimal rounding).
func WithLength(length int32) OptionFunc {
	return func(f *FloatingPoint) {
		f.length = length
	}
}

// WithFormat parses an "Fn(x,y)" format string (e.g. "F9(3,3)") and sets
// length, minimumDecimals, and maximumDecimals in one call.
func WithFormat(format string) OptionFunc {
	return func(f *FloatingPoint) {
		re := regexp.MustCompile(`^F(\d+)\((\d+),(\d+)\)$`)
		m := re.FindStringSubmatch(format)
		if len(m) != 4 {
			return
		}
		if n, err := strconv.ParseInt(m[1], 10, 32); err == nil {
			f.length = int32(n)
		}
		if minD, err := strconv.Atoi(m[2]); err == nil {
			f.minimumDecimals = minD
		}
		if maxD, err := strconv.Atoi(m[3]); err == nil {
			f.maximumDecimals = maxD
		}
	}
}

// Value returns the parsed float64 value.
func (f *FloatingPoint) Value() float64 {
	return f.value
}

func digitsAfterDot(value string) (int, error) {
	var splitMatchCount = 2
	split := strings.Split(value, ".")
	splitLen := len(split)

	// Since Split returns the original value if there's no split, overwrite splitLen
	if splitLen == 1 && split[0] == value {
		splitLen = 0
	}

	switch splitLen {
	case 0:
		return 0, nil
	case splitMatchCount:
		return len(split[1]), nil
	default:
		return 0, errors.New("multiple dots in string")
	}
}

func stringToFloat64(value string, decimals int32) float64 {
	dec, _ := decimal.NewFromString(value)
	dec.Round(decimals)
	result, _ := dec.Float64()
	return result
}

func isValidString(value string) bool {
	reg := regexp.MustCompile("[^0-9.]+")
	result := reg.FindString(value)
	return result == ""
}
