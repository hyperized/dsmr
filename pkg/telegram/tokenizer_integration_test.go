//go:build integration

package telegram_test

import (
	"os"
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizerFromFile(t *testing.T) {
	contents, err := os.ReadFile("../../examples/telegram_v5_0_2.txt")
	require.NoError(t, err)
	tok, err := telegram.NewTokenizer(string(contents))
	require.NoError(t, err)
	assert.NotEmpty(t, tok.Tokens)
}
