package telegram_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/stretchr/testify/assert"
)

func TestFooterNew(t *testing.T) {
	f := telegram.NewFooter("F46A")
	assert.Equal(t, "Footer [CRC: F46A]", f.String())
}

func TestFooterCRC(t *testing.T) {
	f := telegram.NewFooter("ABCD")
	assert.Equal(t, "ABCD", f.CRC())
}
