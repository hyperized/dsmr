package token_test

import (
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram/token"
	"github.com/stretchr/testify/assert"
)

func TestTokenKind(t *testing.T) {
	tok := token.New("/", 0)
	assert.Equal(t, 0, tok.Kind())
}

func TestTokenValue(t *testing.T) {
	tok := token.New("hello", 1)
	assert.Equal(t, "hello", tok.Value())
}
