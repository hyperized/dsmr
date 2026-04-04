package string

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringWithString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		tag      cosem.Cosem
		min      int
		max      int
		expected *String
		fails    bool
	}{
		{
			name:  "S2 tag 9",
			input: "40",
			tag:   cosem.OctetString,
			min:   2,
			max:   2,
			expected: &String{
				tag:           cosem.OctetString,
				minimumLength: 2,
				maximumLength: 2,
				value:         "40",
			},
		},
		{
			name:  "Sn (0..96) tag 9",
			input: "4B384547303034303436333935353037",
			tag:   cosem.OctetString,
			min:   0,
			max:   96,
			expected: &String{
				tag:           cosem.OctetString,
				minimumLength: 0,
				maximumLength: 96,
				value:         "4B384547303034303436333935353037",
			},
		},
		{
			name:  "invalid — exceeds max",
			input: "hello",
			tag:   cosem.OctetString,
			min:   0,
			max:   1,
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(tt.input, WithTag(tt.tag), WithMinimumLength(tt.min), WithMaximumLength(tt.max))
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, s)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, s)
			}
		})
	}
}
