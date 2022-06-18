package telegram

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestTokenizerFromFile(t *testing.T) {
	contents, _ := ioutil.ReadFile("../../examples/telegram_v5_0_2.txt")
	tokens, err := tokenize(string(contents))
	t.Log(tokens)
	if err != nil {
		t.Error(err)
	}
}

func TestTokenizerMissingCharacter(t *testing.T) {
	token, err := tokenize(">")
	t.Log(token)
	if err == nil && len(token.tokens) != 0 {
		t.Error("expected error, but found nil")
	}
}

func TestTokenizerStructure(t *testing.T) {
	tests := []struct {
		input    string
		expected []token
	}{
		{
			input: "/ISk5\\2MT382-1000\n",
			expected: []token{
				{value: "/", kind: Slash},
				{value: "ISk5", kind: Literal},
				{value: "\\", kind: Backslash},
				{value: "2MT382", kind: Literal},
				{value: "-", kind: Dash},
				{value: "1000", kind: Literal},
				{value: "\n", kind: Newline},
			},
		},
	}

	for _, test := range tests {
		t.Logf("Telegram:\n\n%s", test.input)
		actual, err := tokenize(test.input)
		t.Log(test.expected)
		t.Log(actual.tokens)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(test.expected, actual.tokens) {
			t.Errorf("expected %+v, got %+v", test.expected, actual)
		}
	}

}
