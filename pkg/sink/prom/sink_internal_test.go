package prom

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem"
	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMBusDeviceType exercises every branch of the unexported mBusDeviceType helper.
func TestMBusDeviceType(t *testing.T) {
	tg := telegram.NewTelegram()
	assert.Equal(t, "unknown", mBusDeviceType(tg, "1"))

	d, err := telegram.NewData("0-1:24.1.0(INVALID)")
	require.NoError(t, err)
	tg.Add("0-1:24.1.0", d)
	assert.Equal(t, "INVALID", mBusDeviceType(tg, "1"))

	d2, err := telegram.NewData("0-2:24.1.0(003)")
	require.NoError(t, err)
	tg2 := telegram.NewTelegram()
	tg2.Add("0-2:24.1.0", d2)
	assert.Equal(t, "003", mBusDeviceType(tg2, "2"))
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

	v, ok = toFloat64(true)
	assert.False(t, ok)
	assert.InDelta(t, float64(0), v, 0.0001)
}
