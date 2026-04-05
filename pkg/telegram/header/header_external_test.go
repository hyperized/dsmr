package header_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram/header"
	"github.com/stretchr/testify/assert"
)

func TestHeaderNew(t *testing.T) {
	h := header.New(
		header.WithManufacturer("hyperized"),
		header.WithModel("dsmr"),
		header.WithVersion("1.0"),
	)
	assert.Equal(t, "Header [Manufacturer: hyperized, Model: dsmr, Version: 1.0]", h.String())
}

func TestHeaderAccessors(t *testing.T) {
	h := header.New(
		header.WithManufacturer("ISk"),
		header.WithBaudRateID('5'),
		header.WithModel("2MT382-1000"),
		header.WithVersion(""),
	)
	assert.Equal(t, "ISk", h.Manufacturer())
	assert.Equal(t, byte('5'), h.BaudRateID())
	assert.Equal(t, "2MT382-1000", h.Model())
	assert.Empty(t, h.Version())
}
