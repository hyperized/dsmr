package telegram_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"testing/iotest"

	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeSink is a thread-safe Sink for tests. Optional err lets tests exercise
// the parser's per-sink error-logging path without aborting the stream.
type fakeSink struct {
	count atomic.Int32
	err   error
}

func (f *fakeSink) Write(_ *telegram.Telegram) error {
	f.count.Add(1)
	return f.err
}

// TestParseCRCValidation builds a minimal synthetic telegram with CRLF line
// endings, computes its CRC, and verifies that the parser accepts it.
func TestParseCRCValidation(t *testing.T) {
	const le = "\r\n"

	hdr := "/ISk5\\2MT382-1000" + le
	body := le + "1-3:0.2.8(50)" + le
	excl := "!"

	rawForCRC := []byte(hdr + body + excl)
	crcVal := telegram.ComputeCRC(rawForCRC)
	fullTelegram := hdr + body + fmt.Sprintf("!%04X\n", crcVal)

	p := telegram.NewParser(strings.NewReader(fullTelegram),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(true),
	)
	telegrams := slices.Collect(p.Telegrams())
	assert.Len(t, telegrams, 1, "expected exactly one telegram to pass CRC validation")
}

// TestParseEmptyTelegram verifies that a telegram with only a header and footer
// is emitted with an empty data map.
func TestParseEmptyTelegram(t *testing.T) {
	tg := "/ISk5\\2MT382-1000\r\n\r\n!0000\n"

	p := telegram.NewParser(strings.NewReader(tg),
		telegram.WithLineEnding("\r\n"),
		telegram.WithCRCValidation(false),
	)
	telegrams := slices.Collect(p.Telegrams())
	require.Len(t, telegrams, 1, "expected exactly one telegram")
	assert.Equal(t, 0, telegrams[0].Len())
}

// TestParseMalformedLinesSkipped verifies that unparseable lines are silently
// skipped while valid surrounding lines are still extracted.
func TestParseMalformedLinesSkipped(t *testing.T) {
	const le = "\r\n"
	content := "/ISk5\\2MT382-1000" + le +
		le +
		"1-3:0.2.8(50)" + le +
		"not-a-valid-line" + le +
		"9-9:99.99.9(unknown)" + le +
		"!0000\n"

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
	)
	telegrams := slices.Collect(p.Telegrams())
	require.Len(t, telegrams, 1)
	assert.Equal(t, 1, telegrams[0].Len())
	assert.NotNil(t, telegrams[0].Get("1-3:0.2.8"))
}

// TestParseMBusOtherChannels verifies that M-Bus data on channels 2-4 is
// extracted correctly.
func TestParseMBusOtherChannels(t *testing.T) {
	const le = "\r\n"
	content := "/ISk5\\2MT382-1000" + le +
		le +
		"0-2:24.1.0(003)" + le +
		"0-2:24.2.1(101209110000W)(12785.123*m3)" + le +
		"0-3:24.1.0(003)" + le +
		"0-3:24.2.1(101209110100W)(00056.789*m3)" + le +
		"0-4:24.1.0(003)" + le +
		"0-4:24.2.1(101209110200W)(00001.000*m3)" + le +
		"!0000\n"

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
	)
	telegrams := slices.Collect(p.Telegrams())
	require.Len(t, telegrams, 1)

	for _, ch := range []string{"2", "3", "4"} {
		key := "0-" + ch + ":24.2.1"
		d := telegrams[0].Get(key)
		require.NotNil(t, d, "channel %s M-Bus value missing", ch)
		vals := d.Values()
		require.Len(t, vals, 2, "channel %s: expected timestamp + value", ch)
	}
}

// TestParseCRCMismatchDropsTelegram verifies that a telegram with a wrong CRC
// is dropped rather than emitted.
func TestParseCRCMismatchDropsTelegram(t *testing.T) {
	tg := "/ISk5\\2MT382-1000\r\n\r\n1-3:0.2.8(50)\r\n!0000\n"

	p := telegram.NewParser(strings.NewReader(tg),
		telegram.WithLineEnding("\r\n"),
		telegram.WithCRCValidation(true),
	)
	telegrams := slices.Collect(p.Telegrams())
	assert.Empty(t, telegrams, "telegram with wrong CRC should be dropped")
}

// TestParseStreamMultipleSinks verifies that every registered sink receives
// each parsed telegram in registration order.
func TestParseStreamMultipleSinks(t *testing.T) {
	const le = "\r\n"
	content := "/ISk5\\2MT382-1000" + le + le + "!0000\n" +
		"/ISk5\\2MT382-1000" + le + le + "!0000\n"

	a := &fakeSink{}
	b := &fakeSink{}

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
		telegram.WithSink(a),
		telegram.WithSink(b),
	)
	p.ParseStream()

	assert.Equal(t, int32(2), a.count.Load())
	assert.Equal(t, int32(2), b.count.Load())
}

// TestParseStreamSinkErrorLogged verifies that a sink returning an error is
// logged but does not interrupt the stream or starve later sinks.
func TestParseStreamSinkErrorLogged(t *testing.T) {
	const le = "\r\n"
	content := "/ISk5\\2MT382-1000" + le + le + "!0000\n"

	failing := &fakeSink{err: errors.New("boom")}
	healthy := &fakeSink{}

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
		telegram.WithSink(failing),
		telegram.WithSink(healthy),
	)
	p.ParseStream()

	assert.Equal(t, int32(1), failing.count.Load())
	assert.Equal(t, int32(1), healthy.count.Load(), "later sink must still receive the telegram")
}

// TestParseStreamNoSinks verifies that ParseStream is a no-op for the
// telegram body when no sinks are registered.
func TestParseStreamNoSinks(_ *testing.T) {
	const le = "\r\n"
	content := "/ISk5\\2MT382-1000" + le + le + "!0000\n"

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
	)
	p.ParseStream()
}

// TestParseTokenizeErrorOutsideTelegram covers the noise-before-telegram
// branch of the parser's tokenize-error handler. ">" is not a recognised
// tokenizer character, so the line is rejected before any header is seen.
func TestParseTokenizeErrorOutsideTelegram(t *testing.T) {
	const le = "\r\n"
	content := ">>>" + le +
		"/ISk5\\2MT382-1000" + le +
		le +
		"!0000\n"

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
	)
	telegrams := slices.Collect(p.Telegrams())
	assert.Len(t, telegrams, 1)
}

// TestParseTokenizeErrorInsideTelegram covers the tokenize-error-inside-
// telegram branch by injecting an unrecognised character between header and
// footer.
func TestParseTokenizeErrorInsideTelegram(t *testing.T) {
	const le = "\r\n"
	content := "/ISk5\\2MT382-1000" + le +
		le +
		">>>" + le +
		"!0000\n"

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
	)
	telegrams := slices.Collect(p.Telegrams())
	assert.Len(t, telegrams, 1)
}

// TestParseScannerError covers the scanner-error branch (eof with err != nil)
// by feeding the parser a reader that returns a non-EOF error on first read.
func TestParseScannerError(t *testing.T) {
	p := telegram.NewParser(iotest.ErrReader(errors.New("boom")),
		telegram.WithCRCValidation(false),
	)
	telegrams := slices.Collect(p.Telegrams())
	assert.Empty(t, telegrams)
}

// TestTelegramsEarlyBreak exercises the early-return path inside the Telegrams
// iterator when the caller breaks after the first telegram.
func TestTelegramsEarlyBreak(t *testing.T) {
	const le = "\r\n"
	tg1 := "/ISk5\\2MT382-1000" + le + le + "!0000" + le
	tg2 := "/ISk5\\2MT382-1000" + le + le + "!0000" + le
	content := tg1 + tg2

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
	)

	count := 0
	for range p.Telegrams() {
		count++
		break
	}
	assert.Equal(t, 1, count)
}
