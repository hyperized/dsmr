package telegram_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/stretchr/testify/assert"
)

func TestHeaderNew(t *testing.T) {
	h := telegram.NewHeader("hyperized", '5', "dsmr", "1.0")
	assert.Equal(t, "Header [Manufacturer: hyperized, Model: dsmr, Version: 1.0]", h.String())
}

func TestHeaderAccessors(t *testing.T) {
	h := telegram.NewHeader("ISk", '5', "2MT382-1000", "")
	assert.Equal(t, "ISk", h.Manufacturer())
	assert.Equal(t, byte('5'), h.BaudRateID())
	assert.Equal(t, "2MT382-1000", h.Model())
	assert.Empty(t, h.Version())
}
