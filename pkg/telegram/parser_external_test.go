package telegram_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/hyperized/dsmr/pkg/obis"
	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMetrics() *obis.Metrics {
	return obis.Register(prometheus.NewRegistry())
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

// TestParseStreamWithMetrics exercises ParseStream → handleTelegram →
// updateMetrics using a telegram that covers every branch of updateMetrics:
//   - "1-3:0.2.8"          DSMRInfo GaugeVec branch
//   - "0-0:96.1.1"         EquipInfo GaugeVec branch
//   - "0-1:24.1.0"         device-type lookup for MBus label
//   - "0-1:96.1.0"         MBusEquipInfo GaugeVec branch
//   - "0-1:24.2.1"         MBus GaugeVec branch (channel + device_type labels)
//   - "0-1:24.4.0"         MBusValve GaugeVec branch
//   - "1-0:1.8.1"          default branch, FloatingPoint gauge (toFloat64 float path)
//   - "0-0:96.14.0"        default branch, string gauge   (toFloat64 string path)
//   - "0-0:1.0.0"          default branch, no metric name → continue
//   - "1-0:99.97.0"        GenericProfile → TypedValues nil → continue
func TestParseStreamWithMetrics(_ *testing.T) {
	const le = "\r\n"
	content := "/ISk5\\2MT382-1000" + le +
		le +
		"1-3:0.2.8(50)" + le +
		"0-0:96.1.1(4B384547303034303436333935353037)" + le +
		"0-1:24.1.0(003)" + le +
		"0-1:96.1.0(3232323241424344313233343536373839)" + le +
		"0-1:24.2.1(101209112500W)(12785.123*m3)" + le +
		"0-1:24.4.0(1)" + le +
		"1-0:1.8.1(123456.789*kWh)" + le +
		"0-0:96.14.0(0002)" + le +
		"0-0:1.0.0(101209113020W)" + le +
		"1-0:99.97.0(2)(0-0:96.7.19)" + le +
		"!0000\n"

	p := telegram.NewParser(strings.NewReader(content),
		telegram.WithLineEnding(le),
		telegram.WithCRCValidation(false),
		telegram.WithMetrics(newTestMetrics()),
	)
	p.ParseStream()
}

// TestParseStreamNoMetrics exercises the ParseStream nil-metrics branch.
func TestParseStreamNoMetrics(_ *testing.T) {
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
	telegram1 := "/ISk5\\2MT382-1000" + le + le + "!0000" + le
	telegram2 := "/ISk5\\2MT382-1000" + le + le + "!0000" + le
	content := telegram1 + telegram2

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
