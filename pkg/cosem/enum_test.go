package cosem_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnumNew(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  uint8
		fails bool
	}{
		{"zero", "0", 0, false},
		{"max", "255", 255, false},
		{"mid", "42", 42, false},
		{"exceeds uint8", "256", 0, true},
		{"non-numeric", "abc", 0, true},
		{"negative", "-1", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := cosem.NewEnum(tt.input)
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, e)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, e.Value())
		})
	}
}
