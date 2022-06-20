package telegram

import (
	"fmt"
	"unicode"
)

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

type (
	Token struct {
		raw    string
		tokens []token
	}
	token struct {
		value string
		kind  int
	}
)

// TODO: Fix: could not tokenize line: could not identify character '' of type 'uint8'

func tokenize(input string) (Token, error) {
	var ts []token
	length := len(input)
	for counter := 0; counter < length; {
		switch {
		case string(input[counter]) == SlashToken:
			ts = append(ts, token{SlashToken, Slash})
			counter++
		case string(input[counter]) == BackslashToken:
			ts = append(ts, token{BackslashToken, Backslash})
			counter++
		case string(input[counter]) == DashToken:
			ts = append(ts, token{DashToken, Dash})
			counter++
		case string(input[counter]) == ColonToken:
			ts = append(ts, token{ColonToken, Colon})
			counter++
		case string(input[counter]) == ParenthesisOpenToken:
			ts = append(ts, token{ParenthesisOpenToken, ParenthesisOpen})
			counter++
		case string(input[counter]) == ParenthesisCloseToken:
			ts = append(ts, token{ParenthesisCloseToken, ParenthesisClose})
			counter++
		case string(input[counter]) == AsteriskToken:
			ts = append(ts, token{AsteriskToken, Asterisk})
			counter++
		case string(input[counter]) == DotToken:
			ts = append(ts, token{DotToken, Dot})
			counter++
		case string(input[counter]) == ExclamationToken:
			ts = append(ts, token{ExclamationToken, Exclamation})
			counter++
		case rune(input[counter]) == NewlineToken:
			ts = append(ts, token{string(NewlineToken), Newline})
			counter++
		case isLiteral(rune(input[counter])):
			initial := counter
			for counter < length && isLiteral(rune(input[counter])) {
				counter++
			}
			ts = append(ts, token{input[initial:counter], Literal})
		default:
			return Token{}, fmt.Errorf("could not identify character '%c' of type '%T'", input[counter], input[counter])
		}
	}

	return Token{
		raw:    input,
		tokens: ts,
	}, nil
}

func isLiteral(rune rune) bool {
	return unicode.IsLetter(rune) || unicode.IsNumber(rune)
}
