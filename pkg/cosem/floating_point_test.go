package cosem

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsFloatingPointInput(t *testing.T) {
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
			assert.Equal(t, tt.expected, isFloatingPointInput(tt.input))
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

func TestFloatingPointNew(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int32
		min    int
		max    int
		want   float64
		fails  bool
	}{
		{"F9(3,3)", "123456.789", 9, 3, 3, 123456.789, false},
		{"F5(3,3)", "01.193", 5, 3, 3, 1.193, false},
		{"F3(0,0)", "003", 3, 0, 0, 3, false},
		{"F5(0,0)", "00004", 5, 0, 0, 4, false},
		{"F4(1,1)", "220.1", 4, 1, 1, 220.1, false},
		{"decimal mismatch fails", "220.1", 0, 0, 0, 0, true},
		{"multiple dots fails", "1.2.3", 0, 0, 5, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFloatingPoint(tt.input,
				WithLength(tt.length),
				WithMinimumDecimals(tt.min),
				WithMaximumDecimals(tt.max),
			)
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, f)
				return
			}
			require.NoError(t, err)
			assert.InDelta(t, tt.want, f.Value(), 0.0001)
		})
	}
}

func TestFloatingPointWithFormat(t *testing.T) {
	f, err := NewFloatingPoint("123456.789", WithFormat("F9(3,3)"))
	require.NoError(t, err)
	assert.InDelta(t, 123456.789, f.Value(), 0.0001)
}

func TestFloatingPointWithFormatInvalidIgnored(t *testing.T) {
	f, err := NewFloatingPoint("1", WithFormat("INVALID"))
	require.NoError(t, err)
	assert.InDelta(t, 1.0, f.Value(), 0.0001)
}

func TestFloatingPointInvalidCharacters(t *testing.T) {
	f, err := NewFloatingPoint("not-a-number")
	require.Error(t, err)
	assert.Nil(t, f)
}
