package integer_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem/integer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegerValue(t *testing.T) {
	i, err := integer.New("99")
	require.NoError(t, err)
	assert.Equal(t, int64(99), i.Value())
}
