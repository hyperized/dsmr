package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHeaderFromTokenizerEdgeCases exercises branches inside
// headerFromTokenizer that are not reachable through a normal telegram stream.
func TestHeaderFromTokenizerEdgeCases(t *testing.T) {
	// Empty input → no tokens → returns false.
	tok, err := NewTokenizer("")
	require.NoError(t, err)
	h, ok := headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// First token is not a slash → returns false.
	tok, err = NewTokenizer("ABC")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// Slash present but fewer than 2 literals → returns false.
	tok, err = NewTokenizer("/ABC")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// Identifier shorter than 4 characters → returns false.
	tok, err = NewTokenizer("/AB\\model")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.False(t, ok)
	assert.Nil(t, h)

	// Baud-rate ID is not '5' → logs a warning but still returns true.
	tok, err = NewTokenizer("/ISk3\\2MT382-1000")
	require.NoError(t, err)
	h, ok = headerFromTokenizer(tok)
	assert.True(t, ok)
	assert.NotNil(t, h)
}
