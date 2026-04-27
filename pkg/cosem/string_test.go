package cosem_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringNew(t *testing.T) {
	tests := []struct {
		name  string
		input string
		min   int
		max   int
		fails bool
	}{
		{"S2", "40", 2, 2, false},
		{"Sn (0..96)", "4B384547303034303436333935353037", 0, 96, false},
		{"exceeds max", "hello", 0, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := cosem.NewString(tt.input,
				cosem.WithMinimumStringLength(tt.min),
				cosem.WithMaximumStringLength(tt.max),
			)
			if tt.fails {
				require.Error(t, err)
				assert.Nil(t, s)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.input, s.Value())
		})
	}
}
