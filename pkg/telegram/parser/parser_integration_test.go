//go:build integration

package parser_test

import (
	"os"
	"testing"
	"time"

	"github.com/hyperized/dsmr/pkg/telegram/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParse parses the bundled example file with CRC validation disabled
// because the fixture uses LF line endings while DSMR P1 specifies CRLF,
// making stored CRC values not match the computed ones.
func TestParse(t *testing.T) {
	file, err := os.Open("../../../examples/telegram_v5_0_2.txt")
	require.NoError(t, err)
	defer file.Close()

	p := parser.New(file, parser.WithCRCValidation(false))
	telegrams := p.Parse()
	assert.NotEmpty(t, telegrams)
}

// TestParseExampleTelegramValues parses the bundled example file and asserts
// specific OBIS values from the second (complete, poly-phase) telegram.
func TestParseExampleTelegramValues(t *testing.T) {
	file, err := os.Open("../../../examples/telegram_v5_0_2.txt")
	require.NoError(t, err)
	defer file.Close()

	p := parser.New(file, parser.WithCRCValidation(false))
	telegrams := p.Parse()
	require.Len(t, telegrams, 2, "example file must contain exactly 2 telegrams")

	dm := telegrams[1].Data()

	d := dm["1-3:0.2.8"]
	require.NotNil(t, d)
	assert.Equal(t, []string{"50"}, d.Values())

	d = dm["1-0:1.8.1"]
	require.NotNil(t, d)
	assert.Equal(t, []string{"123456.789"}, d.Values())

	d = dm["1-0:1.7.0"]
	require.NotNil(t, d)
	assert.Equal(t, []string{"01.193"}, d.Values())

	d = dm["1-0:32.7.0"]
	require.NotNil(t, d)
	assert.Equal(t, []string{"220.1"}, d.Values())

	d = dm["1-0:51.7.0"]
	require.NotNil(t, d)
	assert.Equal(t, []string{"002"}, d.Values())

	d = dm["1-0:61.7.0"]
	require.NotNil(t, d)
	assert.Equal(t, []string{"03.333"}, d.Values())

	d = dm["0-1:24.2.1"]
	require.NotNil(t, d)
	vals := d.Values()
	require.Len(t, vals, 2)
	assert.Equal(t, "101209112500W", vals[0])
	assert.Equal(t, "12785.123", vals[1])

	d = dm["1-0:99.97.0"]
	require.NotNil(t, d)
	events := d.Events()
	require.Len(t, events, 2)
	assert.Equal(t, 240*time.Second, events[0].Duration)
	assert.Equal(t, 301*time.Second, events[1].Duration)
}
