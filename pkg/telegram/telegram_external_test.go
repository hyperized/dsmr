package telegram_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/hyperized/dsmr/pkg/telegram/data"
	"github.com/hyperized/dsmr/pkg/telegram/footer"
	"github.com/hyperized/dsmr/pkg/telegram/header"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegramString(t *testing.T) {
	h := header.New(header.WithManufacturer("ISk"), header.WithModel("2MT382"), header.WithVersion(""))
	tg := telegram.New(telegram.WithHeader(h))
	assert.Contains(t, tg.String(), "Telegram")
}

func TestTelegramStringWithData(t *testing.T) {
	h := header.New(header.WithManufacturer("ISk"), header.WithModel("2MT382"), header.WithVersion(""))
	tg := telegram.New(telegram.WithHeader(h))
	d, err := data.New("1-3:0.2.8(50)")
	require.NoError(t, err)
	tg.Add("1-3:0.2.8", d)
	s := tg.String()
	assert.Contains(t, s, "Telegram")
	assert.Contains(t, s, "50")
}

func TestTelegramData(t *testing.T) {
	tg := telegram.New()
	d, err := data.New("1-3:0.2.8(50)")
	require.NoError(t, err)
	tg.Add("1-3:0.2.8", d)
	dm := tg.Data()
	assert.Len(t, dm, 1)
	assert.Contains(t, dm, "1-3:0.2.8")
}

func TestTelegramHeader(t *testing.T) {
	h := header.New(header.WithManufacturer("ABC"))
	tg := telegram.New(telegram.WithHeader(h))
	assert.Equal(t, h, tg.Header())
}

func TestTelegramSetFooter(_ *testing.T) {
	f := footer.New(footer.WithCRC("ABCD"))
	tg := telegram.New()
	tg.SetFooter(f)
	// SetFooter is exercised; no public Footer() getter exists.
}
