package cosem_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
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
			maxLength: 4,
			wantBytes: []byte{0x40},
		},
		{
			name:      "valid equipment identifier",
			input:     "4B384547303034303436333935353037",
			maxLength: 16,
			wantBytes: []byte{0x4B, 0x38, 0x45, 0x47, 0x30, 0x30, 0x34, 0x30, 0x34, 0x36, 0x33, 0x39, 0x35, 0x35, 0x30, 0x37},
		},
		{
			name:      "lowercase hex accepted",
			input:     "deadbeef",
			maxLength: 4,
			wantBytes: []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
		{
			name:      "odd-length hex rejected",
			input:     "abc",
			maxLength: 4,
			fails:     true,
		},
		{
			name:      "invalid hex characters",
			input:     "zz",
			maxLength: 4,
			fails:     true,
		},
		{
			name:      "exceeds max length",
			input:     "deadbeef",
			maxLength: 2,
			fails:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := cosem.NewOctetString(tt.input,
				cosem.WithMinimumOctetStringLength(tt.minLength),
				cosem.WithMaximumOctetStringLength(tt.maxLength),
			)
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, o)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantBytes, o.Value())
		})
	}
}

func TestOctetStringString(t *testing.T) {
	o, err := cosem.NewOctetString("deadbeef", cosem.WithMaximumOctetStringLength(4))
	require.NoError(t, err)
	assert.Equal(t, "deadbeef", o.String())
}
