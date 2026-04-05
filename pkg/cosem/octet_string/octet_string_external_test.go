package octet_string_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	octetstring "github.com/hyperized/dsmr/pkg/cosem/octet_string"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOctetStringNew(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		minLength int
		maxLength int
		fails     bool
		wantBytes []byte
	}{
		{
			name:      "valid 2-hex-char (1 byte)",
			input:     "40",
			minLength: 0,
			maxLength: 4,
			wantBytes: []byte{0x40},
		},
		{
			name:      "valid equipment identifier",
			input:     "4B384547303034303436333935353037",
			minLength: 0,
			maxLength: 16,
			wantBytes: []byte{0x4B, 0x38, 0x45, 0x47, 0x30, 0x30, 0x34, 0x30, 0x34, 0x36, 0x33, 0x39, 0x35, 0x35, 0x30, 0x37},
		},
		{
			name:      "lowercase hex accepted",
			input:     "deadbeef",
			minLength: 0,
			maxLength: 4,
			wantBytes: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
		{
			name:      "odd-length hex rejected",
			input:     "abc",
			minLength: 0,
			maxLength: 4,
			fails:     true,
		},
		{
			name:      "invalid hex characters",
			input:     "zz",
			minLength: 0,
			maxLength: 4,
			fails:     true,
		},
		{
			name:      "exceeds max length",
			input:     "deadbeef",
			minLength: 0,
			maxLength: 2,
			fails:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := octetstring.New(tt.input,
				octetstring.WithMinimumLength(tt.minLength),
				octetstring.WithMaximumLength(tt.maxLength),
			)
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, o)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantBytes, o.Value())
			}
		})
	}
}

func TestOctetStringWithTag(t *testing.T) {
	o, err := octetstring.New("ff", octetstring.WithTag(cosem.OctetString), octetstring.WithMaximumLength(4))
	require.NoError(t, err)
	assert.Equal(t, []byte{0xFF}, o.Value())
}

func TestOctetStringString(t *testing.T) {
	o, err := octetstring.New("deadbeef", octetstring.WithMaximumLength(4))
	require.NoError(t, err)
	assert.Equal(t, "deadbeef", o.String())
}
