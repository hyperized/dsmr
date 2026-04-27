package telegram

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizerMissingCharacter(t *testing.T) {
	tok, err := NewTokenizer(">")
	require.Error(t, err)
	assert.Nil(t, tok)
}

func TestTokenizerStructure(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "/ISk5\\2MT382-1000\n",
			expected: []Token{
				{"/", Slash},
				{"ISk5", Literal},
				{"\\", Backslash},
				{"2MT382", Literal},
				{"-", Dash},
				{"1000", Literal},
				{"\n", Newline},
			},
		},
		{
			input: "1-0:1.8.1(123456.789*kWh)",
			expected: []Token{
				{"1", Literal},
				{"-", Dash},
				{"0", Literal},
				{":", Colon},
				{"1", Literal},
				{".", Dot},
				{"8", Literal},
				{".", Dot},
				{"1", Literal},
				{"(", ParenthesisOpen},
				{"123456", Literal},
				{".", Dot},
				{"789", Literal},
				{"*", Asterisk},
				{"kWh", Literal},
				{")", ParenthesisClose},
			},
		},
		{
			input: "!ABCD",
			expected: []Token{
				{"!", Exclamation},
				{"ABCD", Literal},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tok, err := NewTokenizer(tt.input)
			require.NoError(t, err)
			if !reflect.DeepEqual(tt.expected, tok.Tokens) {
				t.Errorf("expected %+v, got %+v", tt.expected, tok.Tokens)
			}
		})
	}
}
