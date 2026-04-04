package enum

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnumNew(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		tag      cosem.Cosem
		expected *Enum
		fails    bool
	}{
		{
			name:  "valid enum value 0",
			input: "0",
			tag:   cosem.Enum,
			expected: &Enum{
				tag:   cosem.Enum,
				value: 0,
			},
		},
		{
			name:  "valid enum value 255",
			input: "255",
			tag:   cosem.Enum,
			expected: &Enum{
				tag:   cosem.Enum,
				value: 255,
			},
		},
		{
			name:  "custom tag",
			input: "42",
			tag:   cosem.Unsigned,
			expected: &Enum{
				tag:   cosem.Unsigned,
				value: 42,
			},
		},
		{
			name:  "exceeds uint8 range",
			input: "256",
			tag:   cosem.Enum,
			fails: true,
		},
		{
			name:  "non-numeric value",
			input: "abc",
			tag:   cosem.Enum,
			fails: true,
		},
		{
			name:  "negative value",
			input: "-1",
			tag:   cosem.Enum,
			fails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := New(tt.input, WithTag(tt.tag))
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, e)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, e)
			}
		})
	}
}
