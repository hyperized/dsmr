package integer

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegerNew(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		tag      cosem.Cosem
		length   int
		expected *Integer
		fails    bool
	}{
		{
			name:   "I4 valid 4-digit integer",
			input:  "0004",
			tag:    cosem.Integer,
			length: 4,
			expected: &Integer{
				tag:    cosem.Integer,
				length: 4,
				value:  4,
			},
		},
		{
			name:   "I2 valid 2-digit integer",
			input:  "50",
			tag:    cosem.Unsigned,
			length: 2,
			expected: &Integer{
				tag:    cosem.Unsigned,
				length: 2,
				value:  50,
			},
		},
		{
			name:   "no length constraint",
			input:  "12345",
			tag:    cosem.Integer,
			length: 0,
			expected: &Integer{
				tag:    cosem.Integer,
				length: 0,
				value:  12345,
			},
		},
		{
			name:   "wrong length",
			input:  "123",
			tag:    cosem.Integer,
			length: 4,
			fails:  true,
		},
		{
			name:   "non-integer value",
			input:  "abc",
			tag:    cosem.Integer,
			length: 0,
			fails:  true,
		},
		{
			name:   "empty value",
			input:  "",
			tag:    cosem.Integer,
			length: 0,
			fails:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := New(tt.input, WithTag(tt.tag), WithLength(tt.length))
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, i)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, i)
			}
		})
	}
}
