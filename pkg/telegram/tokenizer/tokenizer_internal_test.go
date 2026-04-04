package tokenizer

import (
	"reflect"
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizerMissingCharacter(t *testing.T) {
	tok, err := New(">")
	require.Error(t, err)
	assert.Nil(t, tok)
}

func TestTokenizerStructure(t *testing.T) {
	tests := []struct {
		input    string
		expected []*token.Token
	}{
		{
			input: "/ISk5\\2MT382-1000\n",
			expected: []*token.Token{
				token.New("/", Slash),
				token.New("ISk5", Literal),
				token.New("\\", Backslash),
				token.New("2MT382", Literal),
				token.New("-", Dash),
				token.New("1000", Literal),
				token.New("\n", Newline),
			},
		},
		{
			input: "1-0:1.8.1(123456.789*kWh)",
			expected: []*token.Token{
				token.New("1", Literal),
				token.New("-", Dash),
				token.New("0", Literal),
				token.New(":", Colon),
				token.New("1", Literal),
				token.New(".", Dot),
				token.New("8", Literal),
				token.New(".", Dot),
				token.New("1", Literal),
				token.New("(", ParenthesisOpen),
				token.New("123456", Literal),
				token.New(".", Dot),
				token.New("789", Literal),
				token.New("*", Asterisk),
				token.New("kWh", Literal),
				token.New(")", ParenthesisClose),
			},
		},
		{
			input: "!ABCD",
			expected: []*token.Token{
				token.New("!", Exclamation),
				token.New("ABCD", Literal),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tok, err := New(tt.input)
			require.NoError(t, err)
			if !reflect.DeepEqual(tt.expected, tok.tokens) {
				t.Errorf("expected %+v, got %+v", tt.expected, tok.tokens)
			}
		})
	}
}
