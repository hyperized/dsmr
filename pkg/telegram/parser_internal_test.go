package telegram

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHeaderFromTokenizerEdgeCases exercises branches inside
// headerFromTokenizer that are not reachable through a normal telegram stream.
func TestHeaderFromTokenizerEdgeCases(t *testing.T) {
	// Empty input → no tokens → returns false.
	tok, err := NewTokenizer("")
	require.NoError(t, err)
	h, ok := headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// First token is not a slash → returns false.
	tok, err = NewTokenizer("ABC")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// Slash present but fewer than 2 literals → returns false.
	tok, err = NewTokenizer("/ABC")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// Identifier shorter than 4 characters → returns false.
	tok, err = NewTokenizer("/AB\\model")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// Baud-rate ID is not '5' → logs a warning but still returns true.
	tok, err = NewTokenizer("/ISk3\\2MT382-1000")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.True(t, ok)
	assert.NotNil(t, h)
}

// TestMBusDeviceType exercises every branch of the unexported mBusDeviceType helper.
func TestMBusDeviceType(t *testing.T) {
	// Channel absent → "unknown".
	tg := NewTelegram()
	assert.Equal(t, "unknown", mBusDeviceType(tg, "1"))

	// Channel present → raw device-type string is returned regardless of the
	// integer parse outcome (we no longer round-trip through strconv).
	d, err := NewData("0-1:24.1.0(INVALID)")
	require.NoError(t, err)
	tg.Add("0-1:24.1.0", d)
	assert.Equal(t, "INVALID", mBusDeviceType(tg, "1"))

	// Valid device type → original string preserved verbatim.
	d2, err := NewData("0-2:24.1.0(003)")
	require.NoError(t, err)
	tg2 := NewTelegram()
	tg2.Add("0-2:24.1.0", d2)
	assert.Equal(t, "003", mBusDeviceType(tg2, "2"))

	// Channel present but Data has no values → "unknown". Constructed by hand
	// because NewData rejects empty-value lines, so this branch is otherwise
	// unreachable.
	tg3 := NewTelegram()
	tg3.Add("0-3:24.1.0", &Data{})
	assert.Equal(t, "unknown", mBusDeviceType(tg3, "3"))
}

// TestToFloat64 exercises every branch of the unexported toFloat64 helper.
func TestToFloat64(t *testing.T) {
	fp, err := cosem.NewFloatingPoint("1.23",
		cosem.WithMinimumDecimals(2),
		cosem.WithMaximumDecimals(2),
	)
	require.NoError(t, err)
	v, ok := toFloat64(fp)
	assert.True(t, ok)
	assert.InDelta(t, 1.23, v, 0.0001)

	i, err := cosem.NewInteger("42")
	require.NoError(t, err)
	v, ok = toFloat64(i)
	assert.True(t, ok)
	assert.InDelta(t, float64(42), v, 0.0001)

	v, ok = toFloat64("3.14")
	assert.True(t, ok)
	assert.InDelta(t, 3.14, v, 0.0001)

	v, ok = toFloat64("notanumber")
	assert.False(t, ok)
	assert.InDelta(t, float64(0), v, 0.0001)

	v, ok = toFloat64(true) // unhandled type → (0, false)
	assert.False(t, ok)
	assert.InDelta(t, float64(0), v, 0.0001)
}
