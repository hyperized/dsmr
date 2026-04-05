package floating_point_test

import (
	"testing"

	fp "github.com/hyperized/dsmr/pkg/cosem/floating_point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFloatingPointWithFormat(t *testing.T) {
	f, err := fp.New("123456.789", fp.WithFormat("F9(3,3)"))
	require.NoError(t, err)
	assert.InDelta(t, 123456.789, f.Value(), 0.0001)
}

func TestFloatingPointWithFormatInvalidIgnored(t *testing.T) {
	// A malformed format string leaves FloatingPoint with zero constraints,
	// so a whole-number value (0 decimals) is accepted.
	f, err := fp.New("1", fp.WithFormat("INVALID"))
	require.NoError(t, err)
	assert.InDelta(t, 1.0, f.Value(), 0.0001)
}

func TestFloatingPointInvalidCharacters(t *testing.T) {
	// Non-digit, non-dot characters → isValidString returns false → error.
	f, err := fp.New("not-a-number")
	require.Error(t, err)
	assert.Nil(t, f)
}

func TestFloatingPointValue(t *testing.T) {
	f, err := fp.New("3.14",
		fp.WithMinimumDecimals(2),
		fp.WithMaximumDecimals(2),
	)
	require.NoError(t, err)
	assert.InDelta(t, 3.14, f.Value(), 0.0001)
}
