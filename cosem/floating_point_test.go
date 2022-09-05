package cosem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidString(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "no decimals",
			input:    "123456789",
			expected: true,
		},
		{
			name:     "three decimals",
			input:    "123456.789",
			expected: true,
		},
		{
			name:     "no digits",
			input:    "hello world",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := isValidString(test.input)
			assert.Equal(t, test.expected, output)
		})
	}
}

func TestStringToFloat64(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		decimals int32
		expected float64
	}{
		{
			name:     "no decimals",
			input:    "123456789",
			decimals: 0,
			expected: 123456789,
		},
		{
			name:     "three decimals",
			input:    "123456789",
			decimals: 3,
			expected: 123456789.000,
		},
		{
			name:     "no digits",
			input:    "hello world",
			decimals: 3,
			expected: 0.000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := stringToFloat64(test.input, test.decimals)
			assert.Equal(t, test.expected, output)
		})
	}
}

func TestDigitsAfterDots(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected int
		fails    bool
	}{
		{
			name:     "zero",
			input:    "123456789",
			expected: 0,
			fails:    false,
		},
		{
			name:     "one",
			input:    "12345678.9",
			expected: 1,
			fails:    false,
		},
		{
			name:     "more",
			input:    "123456.789",
			expected: 3,
			fails:    false,
		},
		{
			name:     "too much dots",
			input:    "123.456.789",
			expected: 0,
			fails:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := digitsAfterDot(test.input)
			if test.fails {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, output)
		})
	}
}

//nolint:funlen
func TestFloatingPointNewFromString(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		tag      Cosem
		length   int32
		min      int
		max      int
		expected FloatingPoint
		fails    bool
	}{
		{
			name:   "F9(3,3), tag 6",
			input:  "123456.789",
			tag:    DoubleLongUnsigned,
			length: 9,
			min:    3,
			max:    3,
			expected: FloatingPoint{
				Tag:         DoubleLongUnsigned,
				Length:      9,
				MinDecimals: 3,
				MaxDecimals: 3,
				value:       123456.789,
			},
			fails: false,
		},
		{
			name:   "F5(3,3), tag 18",
			input:  "01.193",
			tag:    LongUnsigned,
			length: 5,
			min:    3,
			max:    3,
			expected: FloatingPoint{
				Tag:         LongUnsigned,
				Length:      5,
				MinDecimals: 3,
				MaxDecimals: 3,
				value:       1.193,
			},
			fails: false,
		},
		//
		{
			name:   "F3(0,0), tag 17",
			input:  "003",
			tag:    Unsigned,
			length: 3,
			min:    0,
			max:    0,
			expected: FloatingPoint{
				Tag:         Unsigned,
				Length:      3,
				MinDecimals: 0,
				MaxDecimals: 0,
				value:       3,
			},
			fails: false,
		},
		{
			name:   "F5(0,0), tag 18",
			input:  "00004",
			tag:    LongUnsigned,
			length: 5,
			min:    0,
			max:    0,
			expected: FloatingPoint{
				Tag:         LongUnsigned,
				Length:      5,
				MinDecimals: 0,
				MaxDecimals: 0,
				value:       4,
			},
			fails: false,
		},
		{
			name:   "F4(1,1), tag 18",
			input:  "220.1",
			tag:    LongUnsigned,
			length: 4,
			min:    1,
			max:    1,
			expected: FloatingPoint{
				Tag:         LongUnsigned,
				Length:      4,
				MinDecimals: 1,
				MaxDecimals: 1,
				value:       220.1,
			},
			fails: false,
		},
		{
			name:   "Should fail",
			input:  "220.1",
			tag:    LongUnsigned,
			length: 4,
			min:    0,
			max:    0,
			expected: FloatingPoint{
				Tag:         LongUnsigned,
				Length:      4,
				MinDecimals: 0,
				MaxDecimals: 0,
				value:       0,
			},
			fails: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fp := FloatingPoint{ //nolint:exhaustruct
				Tag:         test.tag,
				Length:      test.length,
				MinDecimals: test.min,
				MaxDecimals: test.max,
			}

			err := fp.WithString(test.input)

			if test.fails {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, fp)
		})
	}
}
