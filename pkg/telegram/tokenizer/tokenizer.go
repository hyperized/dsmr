// Package tokenizer breaks a single DSMR P1 telegram line into a flat slice of
// Tokens.  Each character class (slash, colon, parenthesis, literal alphanumeric
// run, etc.) maps to one of the exported integer constants defined here.
package tokenizer

import (
	"fmt"
	"unicode"

	"github.com/hyperized/dsmr/pkg/telegram/token"
)

// Token kind constants; each maps to a single character class in a telegram line.
const (
	Slash            = iota // /	Characters
	Backslash               // \
	Dash                    // -
	Colon                   // :
	ParenthesisOpen         // (
	ParenthesisClose        // )
	Asterisk                // *
	Dot                     // .
	Exclamation             // !
	Newline                 // Every newline
	Literal                 // Any other set of non-whitespace characters

	SlashToken            = "/"
	BackslashToken        = "\\"
	DashToken             = "-"
	ColonToken            = ":"
	ParenthesisOpenToken  = "("
	ParenthesisCloseToken = ")"
	AsteriskToken         = "*"
	DotToken              = "."
	ExclamationToken      = "!"
	NewlineToken          = '\n'
)

// Tokenizer holds the raw input line and its token slice.
type (
	Tokenizer struct {
		raw    string
		tokens []*token.Token
	}
)

// New tokenizes a single line. It returns an error only when an unrecognised
// character is encountered.
func New(input string) (*Tokenizer, error) {
	t := &Tokenizer{
		raw:    input,
		tokens: []*token.Token{},
	}

	length := len(input)
	for counter := 0; counter < length; {
		b := input[counter]
		switch b {
		case '/':
			t.tokens = append(t.tokens, token.New(SlashToken, Slash))
			counter++
		case '\\':
			t.tokens = append(t.tokens, token.New(BackslashToken, Backslash))
			counter++
		case '-':
			t.tokens = append(t.tokens, token.New(DashToken, Dash))
			counter++
		case ':':
			t.tokens = append(t.tokens, token.New(ColonToken, Colon))
			counter++
		case '(':
			t.tokens = append(t.tokens, token.New(ParenthesisOpenToken, ParenthesisOpen))
			counter++
		case ')':
			t.tokens = append(t.tokens, token.New(ParenthesisCloseToken, ParenthesisClose))
			counter++
		case '*':
			t.tokens = append(t.tokens, token.New(AsteriskToken, Asterisk))
			counter++
		case '.':
			t.tokens = append(t.tokens, token.New(DotToken, Dot))
			counter++
		case '!':
			t.tokens = append(t.tokens, token.New(ExclamationToken, Exclamation))
			counter++
		case '\n':
			t.tokens = append(t.tokens, token.New(string(NewlineToken), Newline))
			counter++
		default:
			if !isLiteral(rune(b)) {
				return nil, fmt.Errorf("tokenizer: could not identify character '%c' of type '%T'", b, b)
			}
			initial := counter
			for counter < length && isLiteral(rune(input[counter])) {
				counter++
			}
			t.tokens = append(t.tokens, token.New(input[initial:counter], Literal))
		}
	}

	return t, nil
}

func isLiteral(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

// Tokens returns the flat slice of tokens parsed from the input line.
func (t *Tokenizer) Tokens() []*token.Token {
	return t.tokens
}

// Raw returns the original unmodified input string passed to New.
func (t *Tokenizer) Raw() string {
	return t.raw
}
