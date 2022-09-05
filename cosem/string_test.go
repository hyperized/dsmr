package cosem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//nolint:funlen
func TestStringWithString(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		tag      Cosem
		min      int
		max      int
		expected String
		fails    bool
	}{
		{
			name:  "S2, tag 9",
			input: "40",
			tag:   OctetString,
			min:   2,
			max:   2,
			expected: String{
				Tag:       OctetString,
				MinLength: 2,
				MaxLength: 2,
				value:     "40",
			},
			fails: false,
		},
		{
			name:  "Sn (n=0..96), tag 9",
			input: "4B384547303034303436333935353037",
			tag:   OctetString,
			min:   0,
			max:   96,
			expected: String{
				Tag:       OctetString,
				MinLength: 0,
				MaxLength: 96,
				value:     "4B384547303034303436333935353037",
			},
			fails: false,
		},
		{
			name:  "Sn (n=0..96), tag 9",
			input: "4B384547303034303436333935353037",
			tag:   OctetString,
			min:   0,
			max:   96,
			expected: String{
				Tag:       OctetString,
				MinLength: 0,
				MaxLength: 96,
				value:     "4B384547303034303436333935353037",
			},
			fails: false,
		},
		{
			name:  "should fail",
			input: "hello",
			tag:   OctetString,
			min:   0,
			max:   1,
			expected: String{
				Tag:       OctetString,
				MinLength: 0,
				MaxLength: 1,
				value:     "",
			},
			fails: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := String{ //nolint:exhaustruct
				Tag:       test.tag,
				MinLength: test.min,
				MaxLength: test.max,
			}

			err := str.WithString(test.input)

			if test.fails {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, str)
		})
	}
}
