package telegram

import (
	"fmt"
	"unicode"
)

// Token kinds produced by the line tokenizer.
const (
	Slash = iota
	Backslash
	Dash
	Colon
	ParenthesisOpen
	ParenthesisClose
	Asterisk
	Dot
	Exclamation
	Newline
	Literal
)

// Token is a single tokenized character or run of characters from a DSMR line.
type Token struct {
	Value string
	Kind  int
}

// Tokenizer holds the original line and the tokens produced from it.
type Tokenizer struct {
	Raw    string
	Tokens []Token
}

// NewTokenizer tokenizes a single DSMR line.
func NewTokenizer(input string) (*Tokenizer, error) {
	t := &Tokenizer{
		Raw:    input,
		Tokens: []Token{},
	}

	length := len(input)
	for counter := 0; counter < length; {
		b := input[counter]
		switch b {
		case '/':
			t.Tokens = append(t.Tokens, Token{"/", Slash})
			counter++
		case '\\':
			t.Tokens = append(t.Tokens, Token{"\\", Backslash})
			counter++
		case '-':
			t.Tokens = append(t.Tokens, Token{"-", Dash})
			counter++
		case ':':
			t.Tokens = append(t.Tokens, Token{":", Colon})
			counter++
		case '(':
			t.Tokens = append(t.Tokens, Token{"(", ParenthesisOpen})
			counter++
		case ')':
			t.Tokens = append(t.Tokens, Token{")", ParenthesisClose})
			counter++
		case '*':
			t.Tokens = append(t.Tokens, Token{"*", Asterisk})
			counter++
		case '.':
			t.Tokens = append(t.Tokens, Token{".", Dot})
			counter++
		case '!':
			t.Tokens = append(t.Tokens, Token{"!", Exclamation})
			counter++
		case '\n':
			t.Tokens = append(t.Tokens, Token{"\n", Newline})
			counter++
		default:
			if !isLiteral(rune(b)) {
				return nil, fmt.Errorf("tokenizer: could not identify character '%c' of type '%T'", b, b)
			}
			initial := counter
			for counter < length && isLiteral(rune(input[counter])) {
				counter++
			}
			t.Tokens = append(t.Tokens, Token{input[initial:counter], Literal})
		}
	}

	return t, nil
}

func isLiteral(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}
