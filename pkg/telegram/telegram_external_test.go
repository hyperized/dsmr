package telegram_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegramString(t *testing.T) {
	tg := telegram.NewTelegram()
	tg.SetHeader(telegram.NewHeader("ISk", '5', "2MT382", ""))
	assert.Contains(t, tg.String(), "Telegram")
}

func TestTelegramStringWithData(t *testing.T) {
	tg := telegram.NewTelegram()
	tg.SetHeader(telegram.NewHeader("ISk", '5', "2MT382", ""))
	d, err := telegram.NewData("1-3:0.2.8(50)")
	require.NoError(t, err)
	tg.Add("1-3:0.2.8", d)
	s := tg.String()
	assert.Contains(t, s, "Telegram")
	assert.Contains(t, s, "50")
}

func TestTelegramData(t *testing.T) {
	tg := telegram.NewTelegram()
	d, err := telegram.NewData("1-3:0.2.8(50)")
	require.NoError(t, err)
	tg.Add("1-3:0.2.8", d)
	assert.Equal(t, 1, tg.Len())
	assert.NotNil(t, tg.Get("1-3:0.2.8"))
	assert.Nil(t, tg.Get("0-0:99.99.9"))
}

func TestTelegramHeader(t *testing.T) {
	h := telegram.NewHeader("ABC", 0, "", "")
	tg := telegram.NewTelegram()
	tg.SetHeader(h)
	assert.Equal(t, h, tg.Header())
}

func TestTelegramSetFooter(_ *testing.T) {
	tg := telegram.NewTelegram()
	tg.SetFooter(telegram.NewFooter("ABCD"))
	// SetFooter is exercised; no public Footer() getter exists.
}
