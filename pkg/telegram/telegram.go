// Package telegram models a complete DSMR P1 telegram: a header, a map of
// parsed OBIS data lines keyed by OBIS identifier, and a footer with the
// CRC-16 checksum.
package telegram

import (
	"strings"

	"github.com/hyperized/dsmr/pkg/telegram/data"
	"github.com/hyperized/dsmr/pkg/telegram/footer"
	"github.com/hyperized/dsmr/pkg/telegram/header"
)

// OptionsFunc configures a Telegram.
type OptionsFunc func(t *Telegram)

// DataMap maps an OBIS identifier string to its parsed Data value.
type DataMap map[string]*data.Data

// Telegram represents a complete, parsed DSMR P1 telegram.
type (
	Telegram struct {
		header *header.Header
		data   DataMap
		footer *footer.Footer
	}
)

// New creates an empty Telegram, applying any provided options.
func New(options ...OptionsFunc) *Telegram {
	t := &Telegram{
		header: nil,
		data:   make(DataMap),
		footer: nil,
	}

	for _, o := range options {
		o(t)
	}

	return t
}

func (t *Telegram) String() string {
	var sb strings.Builder
	sb.WriteString("\nTelegram:\n")
	sb.WriteString(t.header.String())
	for _, d := range t.data {
		sb.WriteByte('\n')
		sb.WriteString(d.String())
	}
	sb.WriteByte('\n')
	return sb.String()
}

// Data returns the map of parsed OBIS data lines keyed by OBIS identifier.
func (t *Telegram) Data() DataMap {
	return t.data
}

// WithHeader sets the header on the telegram.
func WithHeader(h *header.Header) OptionsFunc {
	return func(t *Telegram) {
		t.header = h
	}
}

// WithFooter sets the footer on the telegram.
func WithFooter(f *footer.Footer) OptionsFunc {
	return func(t *Telegram) {
		t.footer = f
	}
}

// SetFooter sets the footer after construction (used by the parser).
func (t *Telegram) SetFooter(f *footer.Footer) {
	t.footer = f
}

// Header returns the parsed identification header.
func (t *Telegram) Header() *header.Header {
	return t.header
}

// Add stores a parsed data line under its OBIS identifier and returns t for chaining.
func (t *Telegram) Add(name string, obis *data.Data) *Telegram {
	t.data[name] = obis
	return t
}
