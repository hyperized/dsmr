// Package parser parses a stream of DSMR P1 telegrams from an io.Reader,
// validating CRC checksums and optionally updating Prometheus metrics.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"slices"
	"strconv"

	"github.com/hyperized/dsmr/pkg/cosem/floating_point"
	"github.com/hyperized/dsmr/pkg/cosem/integer"
	"github.com/hyperized/dsmr/pkg/obis"
	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/hyperized/dsmr/pkg/telegram/crc"
	"github.com/hyperized/dsmr/pkg/telegram/data"
	"github.com/hyperized/dsmr/pkg/telegram/footer"
	"github.com/hyperized/dsmr/pkg/telegram/header"
	"github.com/hyperized/dsmr/pkg/telegram/tokenizer"
)

// Parser reads a stream of DSMR P1 telegrams line by line.
type Parser struct {
	scanner     *bufio.Scanner
	lineEnding  string
	validateCRC bool
	metrics     *obis.Metrics
}

// OptionFunc configures a Parser.
type OptionFunc func(*Parser)

// WithLineEnding sets the line ending used when accumulating raw bytes for CRC
// computation. Defaults to "\r\n" (DSMR spec). Use "\n" for test fixtures
// saved with Unix line endings.
func WithLineEnding(le string) OptionFunc {
	return func(p *Parser) {
		p.lineEnding = le
	}
}

// WithMetrics injects a Metrics instance so parsed telegrams update Prometheus.
func WithMetrics(m *obis.Metrics) OptionFunc {
	return func(p *Parser) {
		p.metrics = m
	}
}

// WithCRCValidation enables or disables CRC validation. Defaults to true.
func WithCRCValidation(enabled bool) OptionFunc {
	return func(p *Parser) {
		p.validateCRC = enabled
	}
}

// New creates a Parser that reads from r.
func New(r io.Reader, opts ...OptionFunc) *Parser {
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

func (p *Parser) line() (*tokenizer.Tokenizer, bool, error) {
	ok := p.scanner.Scan()
	if !ok {
		err := p.scanner.Err()
		return nil, err == nil, err
	}
	t, err := tokenizer.New(p.scanner.Text())
	return t, false, err
}

// headerFromTokenizer parses a DSMR P1 header line per IEC 62056-21.
//
// Format: /MFRb\identification
//   - MFR: 3-character manufacturer code (first 3 chars of the first literal)
//   - b:   baud-rate identifier ('5' in all DSMR P1 implementations)
func headerFromTokenizer(t *tokenizer.Tokenizer) (*header.Header, bool) {
	tokens := t.Tokens()
	if len(tokens) == 0 || tokens[0].Kind() != tokenizer.Slash {
		return nil, false
	}

	var literals []string
	for _, tok := range tokens {
		if tok.Kind() == tokenizer.Literal {
			literals = append(literals, tok.Value())
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

	return header.New(
		header.WithManufacturer(manufacturer),
		header.WithBaudRateID(baudRateID),
		header.WithModel(model),
		header.WithVersion(version),
	), true
}

func dataFromTokenizer(t *tokenizer.Tokenizer) (*data.Data, bool) {
	d, err := data.New(t.Raw())
	if err != nil {
		return nil, false
	}
	return d, true
}

func footerFromTokenizer(t *tokenizer.Tokenizer) (*footer.Footer, bool) {
	tokens := t.Tokens()
	if len(tokens) == 0 || tokens[0].Kind() != tokenizer.Exclamation {
		return nil, false
	}

	var crcVal string
	for _, tok := range tokens {
		if tok.Kind() == tokenizer.Literal {
			crcVal = tok.Value()
		}
	}

	return footer.New(footer.WithCRC(crcVal)), true
}

// Telegrams returns an iterator over every complete telegram in the stream.
// The iterator starts the internal parser goroutine on first iteration and
// stops when the underlying reader reaches EOF or the caller breaks early.
func (p *Parser) Telegrams() iter.Seq[*telegram.Telegram] {
	return func(yield func(*telegram.Telegram) bool) {
		ch := make(chan *telegram.Telegram)
		go p.parseLines(ch)
		for t := range ch {
			if !yield(t) {
				return
			}
		}
	}
}

// Parse collects all telegrams from the stream and returns them.
func (p *Parser) Parse() []*telegram.Telegram {
	return slices.Collect(p.Telegrams())
}

// ParseStream continuously processes telegrams, updating Prometheus metrics
// for each complete telegram.
func (p *Parser) ParseStream() {
	for t := range p.Telegrams() {
		go p.handleTelegram(t)
	}
}

func (p *Parser) handleTelegram(t *telegram.Telegram) {
	slog.Debug("telegram parsed", "content", t.String())
	if p.metrics != nil {
		p.updateMetrics(t)
	}
}

// updateMetrics iterates the telegram's DataMap and updates all registered
// Prometheus metrics.
func (p *Parser) updateMetrics(t *telegram.Telegram) {
	for id, d := range t.Data() {
		tv := d.TypedValues()
		if len(tv) == 0 {
			continue
		}

		switch id {
		case "1-3:0.2.8":
			if v, ok := tv[0].(string); ok {
				p.metrics.DSMRInfo().WithLabelValues(v).Set(1)
			}
		case "0-0:96.1.1":
			if v, ok := tv[0].(string); ok {
				p.metrics.EquipInfo().WithLabelValues(v).Set(1)
			}
		case "0-1:24.2.1", "0-2:24.2.1", "0-3:24.2.1", "0-4:24.2.1":
			if len(tv) >= 2 {
				if vs, ok := tv[1].(string); ok {
					if f, err := strconv.ParseFloat(vs, 64); err == nil {
						ch := id[2:3]
						deviceType := mBusDeviceType(t, ch)
						p.metrics.MBus().WithLabelValues(ch, deviceType).Set(f)
					}
				}
			}
		case "0-1:24.4.0", "0-2:24.4.0", "0-3:24.4.0", "0-4:24.4.0":
			if f, ok := toFloat64(tv[0]); ok {
				p.metrics.MBusValve().WithLabelValues(id[2:3]).Set(f)
			}
		default:
			metricName := d.MetricName()
			if metricName == "" {
				continue
			}
			if gauge, ok := p.metrics.Gauges()[metricName]; ok {
				if f, ok := toFloat64(tv[0]); ok {
					gauge.Set(f)
				}
			}
		}
	}
}

const deviceTypeUnknown = "unknown"

// mBusDeviceType looks up the device-type value for the given M-Bus channel
// (e.g. "1") from the telegram's DataMap and returns it as a zero-padded
// 3-digit string (e.g. "003" for gas). Falls back to deviceTypeUnknown if the
// device-type line is absent or could not be parsed.
func mBusDeviceType(t *telegram.Telegram, ch string) string {
	obisID := "0-" + ch + ":24.1.0"
	d, ok := t.Data()[obisID]
	if !ok {
		return deviceTypeUnknown
	}
	tv := d.TypedValues()
	if len(tv) == 0 {
		return deviceTypeUnknown
	}
	if i, ok := tv[0].(*integer.Integer); ok {
		return fmt.Sprintf("%03d", i.Value())
	}
	return deviceTypeUnknown
}

func toFloat64(v any) (float64, bool) {
	switch x := v.(type) {
	case *floating_point.FloatingPoint:
		return x.Value(), true
	case *integer.Integer:
		return float64(x.Value()), true
	case string:
		if f, err := strconv.ParseFloat(x, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func (p *Parser) parseLines(ch chan *telegram.Telegram) {
	defer close(ch)

	var current *telegram.Telegram
	var rawBuf []byte

	for {
		t, eof, err := p.line()

		if eof {
			slog.Debug("EOF reached")
			return
		}

		if err != nil {
			slog.Error("scan error", "err", err)
			continue
		}

		// A header line starts a new telegram and resets accumulation.
		if h, ok := headerFromTokenizer(t); ok {
			current = telegram.New(telegram.WithHeader(h))
			rawBuf = []byte(t.Raw() + p.lineEnding)
			continue
		}

		if current == nil {
			continue
		}

		// Footer line: finalise CRC, validate, emit.
		if f, ok := footerFromTokenizer(t); ok {
			rawBuf = append(rawBuf, '!')
			if p.validateCRC && !crc.Valid(rawBuf, f.CRC()) {
				slog.Warn("CRC mismatch, telegram dropped",
					"expected", f.CRC(),
					"computed", fmt.Sprintf("%04X", crc.Compute(rawBuf)))
				current = nil
				rawBuf = nil
				continue
			}
			current.SetFooter(f)
			ch <- current
			current = nil
			rawBuf = nil
			continue
		}

		// Empty or data line — accumulate bytes for CRC.
		rawBuf = append(rawBuf, []byte(t.Raw()+p.lineEnding)...)

		if len(t.Tokens()) == 0 {
			continue
		}

		if d, ok := dataFromTokenizer(t); ok {
			current.Add(d.Identifier(), d)
		}
	}
}
