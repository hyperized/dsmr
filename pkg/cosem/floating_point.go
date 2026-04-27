package cosem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// FloatingPointOption configures a FloatingPoint value.
type FloatingPointOption func(f *FloatingPoint)

// FloatingPoint is a parsed COSEM floating-point value with format constraints.
type FloatingPoint struct {
	length          int32
	minimumDecimals int
	maximumDecimals int
	value           float64
}

// floatingPointFormatRE matches "Fn(x,y)" format strings.
var floatingPointFormatRE = regexp.MustCompile(`^F(\d+)\((\d+),(\d+)\)$`)

// NewFloatingPoint parses value as a decimal floating-point number. The value
// must contain only digits and at most one dot. Decimal-digit count must
// satisfy any min/max constraints set via options.
func NewFloatingPoint(value string, options ...FloatingPointOption) (*FloatingPoint, error) {
	f := &FloatingPoint{}
	for _, o := range options {
		o(f)
	}

	if !isFloatingPointInput(value) {
		return nil, errors.New("the provided string can only contain digits and dots")
	}
	digits, err := digitsAfterDot(value)
	if err != nil {
		return nil, err
	}
	if digits < f.minimumDecimals || digits > f.maximumDecimals {
		return nil, fmt.Errorf("decimal point incorrect, expected to be between %d and %d, found %d",
			f.minimumDecimals, f.maximumDecimals, digits)
	}

	f.value = stringToFloat64(value, f.length)
	return f, nil
}

// WithMinimumDecimals sets the minimum number of digits after the decimal point.
func WithMinimumDecimals(minimum int) FloatingPointOption {
	return func(f *FloatingPoint) { f.minimumDecimals = minimum }
}

// WithMaximumDecimals sets the maximum number of digits after the decimal point.
func WithMaximumDecimals(maximum int) FloatingPointOption {
	return func(f *FloatingPoint) { f.maximumDecimals = maximum }
}

// WithLength sets the total field width (used for decimal rounding).
func WithLength(length int32) FloatingPointOption {
	return func(f *FloatingPoint) { f.length = length }
}

// WithFormat parses an "Fn(x,y)" format string (e.g. "F9(3,3)") and sets
// length, minimumDecimals, and maximumDecimals in one call.
func WithFormat(format string) FloatingPointOption {
	return func(f *FloatingPoint) {
		m := floatingPointFormatRE.FindStringSubmatch(format)
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
	const splitMatchCount = 2
	split := strings.Split(value, ".")
	splitLen := len(split)
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

var floatingPointInputRE = regexp.MustCompile("[^0-9.]+")

func isFloatingPointInput(value string) bool {
	return floatingPointInputRE.FindString(value) == ""
}
