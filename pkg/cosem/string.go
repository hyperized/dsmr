package cosem

import "fmt"

// StringOption configures a String.
type StringOption func(s *String)

// String is a parsed COSEM VisibleString value with optional length constraints.
type String struct {
	minimumLength int
	maximumLength int
	value         string
}

// NewString creates a String from value, applying length constraints if set.
func NewString(value string, options ...StringOption) (*String, error) {
	s := &String{value: value}
	for _, o := range options {
		o(s)
	}

	strLen := len(value)
	if strLen < s.minimumLength || strLen > s.maximumLength {
		return nil, fmt.Errorf("character count exceeds range, expected to be between %d and %d, found %d",
			s.minimumLength, s.maximumLength, strLen)
	}
	return s, nil
}

// WithMinimumStringLength sets the minimum character count (inclusive).
func WithMinimumStringLength(length int) StringOption {
	return func(s *String) { s.minimumLength = length }
}

// WithMaximumStringLength sets the maximum character count (inclusive).
func WithMaximumStringLength(length int) StringOption {
	return func(s *String) { s.maximumLength = length }
}

// Value returns the underlying string value.
func (s *String) Value() string { return s.value }
