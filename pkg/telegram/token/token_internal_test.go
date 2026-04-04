package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenNew(t *testing.T) {
	tok := New("hello", 0)
	assert.Equal(t, &Token{value: "hello", kind: 0}, tok)
}
