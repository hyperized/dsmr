package cosem_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegerNew(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		want   int64
		fails  bool
	}{
		{"I4 valid 4-digit integer", "0004", 4, 4, false},
		{"I2 valid 2-digit integer", "50", 2, 50, false},
		{"no length constraint", "12345", 0, 12345, false},
		{"wrong length", "123", 4, 0, true},
		{"non-integer value", "abc", 0, 0, true},
		{"empty value", "", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := cosem.NewInteger(tt.input, cosem.WithIntegerLength(tt.length))
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, i)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, i.Value())
		})
	}
}
