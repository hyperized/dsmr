package telegram_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizerTokens(t *testing.T) {
	tok, err := telegram.NewTokenizer("/ISk5\\2MT382-1000")
	require.NoError(t, err)
	assert.NotEmpty(t, tok.Tokens)
}

func TestTokenizerRaw(t *testing.T) {
	const input = "1-0:1.8.1(123456.789*kWh)"
	tok, err := telegram.NewTokenizer(input)
	require.NoError(t, err)
	assert.Equal(t, input, tok.Raw)
}

func TestTokenizerEmpty(t *testing.T) {
	tok, err := telegram.NewTokenizer("")
	require.NoError(t, err)
	assert.Empty(t, tok.Tokens)
}
