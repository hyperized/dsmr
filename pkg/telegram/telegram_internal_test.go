package telegram

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram/data"
	"github.com/hyperized/dsmr/pkg/telegram/footer"
	"github.com/hyperized/dsmr/pkg/telegram/header"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegramNew(t *testing.T) {
	h := header.New(
		header.WithManufacturer("test"),
		header.WithModel("demo"),
		header.WithVersion("1.0"),
	)
	f := footer.New(footer.WithCRC("B33F"))

	one, err := data.New("0-0:96.13.0(303132333435363738393A3B3C3D3E3F)")
	require.NoError(t, err)
	two, err := data.New("0-0:96.7.21(00004)")
	require.NoError(t, err)

	tg := New(WithHeader(h), WithFooter(f)).Add("0-0:96.13.0", one).Add("0-0:96.7.21", two)

	assert.Equal(t, &Telegram{
		header: h,
		data: DataMap{
			"0-0:96.13.0": one,
			"0-0:96.7.21": two,
		},
		footer: f,
	}, tg)
}
