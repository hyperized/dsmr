package enum_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnumValue(t *testing.T) {
	e, err := enum.New("42")
	require.NoError(t, err)
	assert.Equal(t, uint8(42), e.Value())
}
