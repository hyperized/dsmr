package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegramAddPreservesInsertionOrder(t *testing.T) {
	h := NewHeader("test", '5', "demo", "1.0")
	f := NewFooter("B33F")

	one, err := NewData("0-0:96.13.0(303132333435363738393A3B3C3D3E3F)")
	require.NoError(t, err)
	two, err := NewData("0-0:96.7.21(00004)")
	require.NoError(t, err)

	tg := NewTelegram()
	tg.SetHeader(h)
	tg.SetFooter(f)
	tg.Add("0-0:96.13.0", one)
	tg.Add("0-0:96.7.21", two)

	assert.Equal(t, h, tg.Header())
	assert.Equal(t, 2, tg.Len())
	assert.Same(t, one, tg.Get("0-0:96.13.0"))
	assert.Same(t, two, tg.Get("0-0:96.7.21"))

	var ids []string
	for id := range tg.All() {
		ids = append(ids, id)
	}
	assert.Equal(t, []string{"0-0:96.13.0", "0-0:96.7.21"}, ids)
}

func TestTelegramAllEarlyBreak(t *testing.T) {
	d1, err := NewData("1-3:0.2.8(50)")
	require.NoError(t, err)
	d2, err := NewData("0-0:96.7.21(00004)")
	require.NoError(t, err)

	tg := NewTelegram()
	tg.Add("1-3:0.2.8", d1)
	tg.Add("0-0:96.7.21", d2)

	count := 0
	for range tg.All() {
		count++
		break
	}
	assert.Equal(t, 1, count)
}

func TestTelegramAddOverwriteKeepsPosition(t *testing.T) {
	d1, err := NewData("1-3:0.2.8(50)")
	require.NoError(t, err)
	d2, err := NewData("0-0:96.7.21(00004)")
	require.NoError(t, err)
	d1Replacement, err := NewData("1-3:0.2.8(50)")
	require.NoError(t, err)

	tg := NewTelegram()
	tg.Add("1-3:0.2.8", d1)
	tg.Add("0-0:96.7.21", d2)
	tg.Add("1-3:0.2.8", d1Replacement)

	assert.Same(t, d1Replacement, tg.Get("1-3:0.2.8"))
	assert.Equal(t, 2, tg.Len())

	var ids []string
	for id := range tg.All() {
		ids = append(ids, id)
	}
	assert.Equal(t, []string{"1-3:0.2.8", "0-0:96.7.21"}, ids)
}
