// Package token defines the Token type produced by the tokenizer.
// Each token carries a kind (an integer constant defined in the tokenizer
// package) and the literal string value it was scanned from.
package token

// Token is a single lexical unit from a DSMR telegram line.
type Token struct {
	value string
	kind  int
}

// New creates a Token with the given value and kind constant.
func New(value string, kind int) *Token {
	t := &Token{
		value: value,
		kind:  kind,
	}

	return t
}

// Kind returns the integer kind constant (one of the tokenizer package constants).
func (t *Token) Kind() int {
	return t.kind
}

// Value returns the literal string scanned for this token.
func (t *Token) Value() string {
	return t.value
}
