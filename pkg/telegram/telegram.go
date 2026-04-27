// Package telegram models a complete DSMR P1 telegram and provides the parser,
// tokenizer, CRC validator, and supporting types used to read telegrams from
// an io.Reader and produce parsed values.
package telegram

import (
	"iter"
	"strings"
)

// Telegram represents a complete, parsed DSMR P1 telegram: a header, an
// insertion-ordered set of parsed OBIS data lines keyed by OBIS identifier,
// and a footer with the CRC-16 checksum.
type Telegram struct {
	header *Header
	keys   []string
	data   map[string]*Data
	footer *Footer
}

// NewTelegram creates an empty Telegram. Header, data lines and footer are
// added incrementally by the parser.
func NewTelegram() *Telegram {
	return &Telegram{data: map[string]*Data{}}
}

func (t *Telegram) String() string {
	var sb strings.Builder
	sb.WriteString("\nTelegram:\n")
	if t.header != nil {
		sb.WriteString(t.header.String())
	}
	for _, k := range t.keys {
		sb.WriteByte('\n')
		sb.WriteString(t.data[k].String())
	}
	sb.WriteByte('\n')
	return sb.String()
}

// Get returns the data line for an OBIS identifier, or nil when absent.
func (t *Telegram) Get(id string) *Data {
	return t.data[id]
}

// All yields every (identifier, data) pair in insertion order.
func (t *Telegram) All() iter.Seq2[string, *Data] {
	return func(yield func(string, *Data) bool) {
		for _, k := range t.keys {
			if !yield(k, t.data[k]) {
				return
			}
		}
	}
}

// Len returns the number of parsed data lines in this telegram.
func (t *Telegram) Len() int { return len(t.keys) }

// SetHeader sets the header on the telegram.
func (t *Telegram) SetHeader(h *Header) { t.header = h }

// SetFooter sets the footer on the telegram.
func (t *Telegram) SetFooter(f *Footer) { t.footer = f }

// Header returns the parsed identification header.
func (t *Telegram) Header() *Header { return t.header }

// Add stores a parsed data line under its OBIS identifier and returns t for chaining.
// Re-adding an existing identifier overwrites the value but preserves its position.
func (t *Telegram) Add(id string, d *Data) *Telegram {
	if _, exists := t.data[id]; !exists {
		t.keys = append(t.keys, id)
	}
	t.data[id] = d
	return t
}
