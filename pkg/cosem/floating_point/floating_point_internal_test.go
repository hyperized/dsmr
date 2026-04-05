package floating_point

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsValidString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"no decimals", "123456789", true},
		{"three decimals", "123456.789", true},
		{"no digits", "hello world", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isValidString(tt.input))
		})
	}
}

func TestStringToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		decimals int32
		expected float64
	}{
		{"no decimals", "123456789", 0, 123456789},
		{"three decimals", "123456789", 3, 123456789.000},
		{"no digits", "hello world", 3, 0.000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.expected, stringToFloat64(tt.input, tt.decimals), 0.0001)
		})
	}
}

func TestDigitsAfterDot(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		fails    bool
	}{
		{"zero", "123456789", 0, false},
		{"one", "12345678.9", 1, false},
		{"more", "123456.789", 3, false},
		{"too many dots", "123.456.789", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := digitsAfterDot(tt.input)
			if tt.fails {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expected, out)
		})
	}
}

//nolint:funlen
func TestFloatingPointNewFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		tag      cosem.Cosem
		length   int32
		min      int
		max      int
		expected *FloatingPoint
		fails    bool
	}{
		{
			name:   "F9(3,3) tag 6",
			input:  "123456.789",
			tag:    cosem.DoubleLongUnsigned,
			length: 9,
			min:    3,
			max:    3,
			expected: &FloatingPoint{
				tag:             cosem.DoubleLongUnsigned,
				length:          9,
				minimumDecimals: 3,
				maximumDecimals: 3,
				value:           123456.789,
			},
		},
		{
			name:   "F5(3,3) tag 18",
			input:  "01.193",
			tag:    cosem.LongUnsigned,
			length: 5,
			min:    3,
			max:    3,
			expected: &FloatingPoint{
				tag:             cosem.LongUnsigned,
				length:          5,
				minimumDecimals: 3,
				maximumDecimals: 3,
				value:           1.193,
			},
		},
		{
			name:   "F3(0,0) tag 17",
			input:  "003",
			tag:    cosem.Unsigned,
			length: 3,
			min:    0,
			max:    0,
			expected: &FloatingPoint{
				tag:             cosem.Unsigned,
				length:          3,
				minimumDecimals: 0,
				maximumDecimals: 0,
				value:           3,
			},
		},
		{
			name:   "F5(0,0) tag 18",
			input:  "00004",
			tag:    cosem.LongUnsigned,
			length: 5,
			min:    0,
			max:    0,
			expected: &FloatingPoint{
				tag:             cosem.LongUnsigned,
				length:          5,
				minimumDecimals: 0,
				maximumDecimals: 0,
				value:           4,
			},
		},
		{
			name:   "F4(1,1) tag 18",
			input:  "220.1",
			tag:    cosem.LongUnsigned,
			length: 4,
			min:    1,
			max:    1,
			expected: &FloatingPoint{
				tag:             cosem.LongUnsigned,
				length:          4,
				minimumDecimals: 1,
				maximumDecimals: 1,
				value:           220.1,
			},
		},
		{
			name:  "decimal count mismatch fails",
			input: "220.1",
			tag:   cosem.LongUnsigned,
			fails: true,
		},
		{
			name:  "multiple dots in value fails",
			input: "1.2.3",
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(tt.input,
				WithTag(tt.tag),
				WithLength(tt.length),
				WithMinimumDecimals(tt.min),
				WithMaximumDecimals(tt.max),
			)
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, f)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, f)
			}
		})
	}
}
