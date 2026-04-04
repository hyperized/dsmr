package parser

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/cosem/floating_point"
	"github.com/hyperized/dsmr/pkg/cosem/integer"
	"github.com/hyperized/dsmr/pkg/telegram/tokenizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHeaderFromTokenizerEdgeCases exercises branches inside
// headerFromTokenizer that are not reachable through a normal telegram stream.
func TestHeaderFromTokenizerEdgeCases(t *testing.T) {
	// Empty input → no tokens → returns false.
	tok, err := tokenizer.New("")
	require.NoError(t, err)
	h, ok := headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// First token is not a slash → returns false.
	tok, err = tokenizer.New("ABC")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// Slash present but fewer than 2 literals → returns false.
	tok, err = tokenizer.New("/ABC")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// Identifier shorter than 4 characters → returns false.
	tok, err = tokenizer.New("/AB\\model")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// Baud-rate ID is not '5' → logs a warning but still returns true.
	tok, err = tokenizer.New("/ISk3\\2MT382-1000")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.True(t, ok)
	assert.NotNil(t, h)
}

// TestToFloat64 exercises every branch of the unexported toFloat64 helper.
func TestToFloat64(t *testing.T) {
	fp, err := floating_point.New("1.23",
		floating_point.WithMinimumDecimals(2),
		floating_point.WithMaximumDecimals(2),
	)
	require.NoError(t, err)
	v, ok := toFloat64(fp)
	assert.True(t, ok)
	assert.InDelta(t, 1.23, v, 0.0001)

	i, err := integer.New("42")
	require.NoError(t, err)
	v, ok = toFloat64(i)
	assert.True(t, ok)
	assert.InDelta(t, float64(42), v, 0.0001)

	v, ok = toFloat64("3.14")
	assert.True(t, ok)
	assert.InDelta(t, 3.14, v, 0.0001)

	v, ok = toFloat64("notanumber")
	assert.False(t, ok)
	assert.InDelta(t, float64(0), v, 0.0001)

	v, ok = toFloat64(true) // unhandled type → (0, false)
	assert.False(t, ok)
	assert.InDelta(t, float64(0), v, 0.0001)
}
