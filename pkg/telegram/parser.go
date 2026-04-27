package telegram

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"log/slog"
)

// Parser reads a stream of DSMR P1 telegrams line by line.
type Parser struct {
	scanner     *bufio.Scanner
	lineEnding  string
	validateCRC bool
	sinks       []Sink
}

// OptionFunc configures a Parser.
type OptionFunc func(*Parser)

// WithLineEnding sets the line ending used when accumulating raw bytes for CRC
// computation. Defaults to "\r\n" (DSMR spec). Use "\n" for test fixtures
// saved with Unix line endings.
func WithLineEnding(le string) OptionFunc {
	return func(p *Parser) { p.lineEnding = le }
}

// WithSink registers a Sink to receive every successfully parsed telegram.
// May be called multiple times; sinks are invoked in registration order and
// per-sink errors are logged but do not interrupt the stream.
func WithSink(s Sink) OptionFunc {
	return func(p *Parser) { p.sinks = append(p.sinks, s) }
}

// WithCRCValidation enables or disables CRC validation. Defaults to true.
func WithCRCValidation(enabled bool) OptionFunc {
	return func(p *Parser) { p.validateCRC = enabled }
}

// NewParser creates a Parser that reads from r.
func NewParser(r io.Reader, opts ...OptionFunc) *Parser {
	p := &Parser{
		scanner:     bufio.NewScanner(r),
		lineEnding:  "\r\n",
		validateCRC: true,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// line returns the next tokenizer for a scanned line, or eof=true when the
// underlying scanner is exhausted (with err set if the scanner failed).
func (p *Parser) line() (tok *Tokenizer, eof bool, err error) {
	if !p.scanner.Scan() {
		return nil, true, p.scanner.Err()
	}
	tok, err = NewTokenizer(p.scanner.Text())
	return tok, false, err
}

// headerFromTokenizer parses a DSMR P1 header line per IEC 62056-21.
//
// Format: /MFRb\identification
//   - MFR: 3-character manufacturer code (first 3 chars of the first literal)
//   - b:   baud-rate identifier ('5' in all DSMR P1 implementations)
func headerFromTokenizer(t *Tokenizer) (*Header, bool) {
	tokens := t.Tokens
	if len(tokens) == 0 || tokens[0].Kind != Slash {
		return nil, false
	}

	var literals []string
	for _, tok := range tokens {
		if tok.Kind == Literal {
			literals = append(literals, tok.Value)
		}
	}

	if len(literals) < 2 {
		return nil, false
	}

	ident := literals[0]
	if len(ident) < 4 {
		return nil, false
	}

	manufacturer := ident[:3]
	baudRateID := ident[3]

	if baudRateID != '5' {
		slog.Warn("unexpected baud rate identifier", "got", string(baudRateID), "expected", "5")
	}

	model := literals[1]
	var version string
	if len(literals) >= 3 {
		version = literals[2]
	}

	return NewHeader(manufacturer, baudRateID, model, version), true
}

func dataFromTokenizer(t *Tokenizer) (*Data, bool) {
	d, err := NewData(t.Raw)
	if err != nil {
		return nil, false
	}
	return d, true
}

func footerFromTokenizer(t *Tokenizer) (*Footer, bool) {
	tokens := t.Tokens
	if len(tokens) == 0 || tokens[0].Kind != Exclamation {
		return nil, false
	}

	var crcVal string
	for _, tok := range tokens {
		if tok.Kind == Literal {
			crcVal = tok.Value
		}
	}

	return NewFooter(crcVal), true
}

// Telegrams returns an iterator over every complete telegram in the stream.
// Iteration is synchronous: scanning happens inside the yield loop, so an
// early break stops scanning immediately and leaks no goroutines.
func (p *Parser) Telegrams() iter.Seq[*Telegram] {
	return func(yield func(*Telegram) bool) {
		var current *Telegram
		var rawBuf []byte

		for {
			t, eof, err := p.line()
			if eof {
				if err != nil {
					slog.Error("scan error", "err", err)
				}
				return
			}
			if err != nil {
				if current == nil {
					slog.Debug("ignoring noise before telegram start", "err", err)
				} else {
					slog.Debug("tokenize error inside telegram", "err", err)
				}
				continue
			}

			if h, ok := headerFromTokenizer(t); ok {
				current = NewTelegram()
				current.SetHeader(h)
				rawBuf = []byte(t.Raw + p.lineEnding)
				continue
			}

			if current == nil {
				continue
			}

			if f, ok := footerFromTokenizer(t); ok {
				rawBuf = append(rawBuf, '!')
				if p.validateCRC && !ValidCRC(rawBuf, f.CRC()) {
					slog.Warn("CRC mismatch, telegram dropped",
						"expected", f.CRC(),
						"computed", fmt.Sprintf("%04X", ComputeCRC(rawBuf)))
					current = nil
					rawBuf = nil
					continue
				}
				current.SetFooter(f)
				if !yield(current) {
					return
				}
				current = nil
				rawBuf = nil
				continue
			}

			rawBuf = append(rawBuf, []byte(t.Raw+p.lineEnding)...)
			if len(t.Tokens) == 0 {
				continue
			}
			if d, ok := dataFromTokenizer(t); ok {
				current.Add(d.Identifier(), d)
			}
		}
	}
}

// ParseStream continuously processes telegrams, dispatching each one to every
// registered sink in order. Per-sink errors are logged and do not interrupt
// the stream.
func (p *Parser) ParseStream() {
	for t := range p.Telegrams() {
		slog.Debug("telegram parsed", "content", t.String())
		for _, s := range p.sinks {
			if err := s.Write(t); err != nil {
				slog.Warn("sink write failed", "err", err)
			}
		}
	}
}
