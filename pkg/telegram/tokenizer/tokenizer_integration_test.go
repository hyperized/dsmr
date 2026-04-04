//go:build integration

package tokenizer_test

import (
	"os"
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram/tokenizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizerFromFile(t *testing.T) {
	contents, err := os.ReadFile("../../examples/telegram_v5_0_2.txt")
	require.NoError(t, err)
	tok, err := tokenizer.New(string(contents))
	require.NoError(t, err)
	assert.NotEmpty(t, tok.Tokens())
}
