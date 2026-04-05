// Package string implements the COSEM VisibleString type.
// Values are plain ASCII strings whose character count can be constrained to a
// [min, max] range.
package string

import (
	"fmt"

	"github.com/hyperized/dsmr/pkg/cosem"
)

// OptionFunc configures a String.
type OptionFunc func(s *String)

// String holds a COSEM VisibleString value with optional length constraints.
type String struct {
	tag           cosem.Cosem
	minimumLength int
	maximumLength int
	value         string
}

// New creates a String from value, applying length constraints if set.
func New(value string, options ...OptionFunc) (*String, error) {
	s := &String{
		tag:           cosem.VisibleString,
		minimumLength: 0,
		maximumLength: 0,
		value:         value,
	}

	for _, o := range options {
		o(s)
	}

	strLen := len(value)
	if strLen < s.minimumLength || strLen > s.maximumLength {
		return nil, fmt.Errorf("character count exceeds range, expected to be between %d and %d, found %d",
			s.minimumLength, s.maximumLength, strLen,
		)
	}

	return s, nil
}

// WithTag overrides the default COSEM tag for this string.
func WithTag(tag cosem.Cosem) OptionFunc {
	return func(s *String) {
		s.tag = tag
	}
}

// WithMinimumLength sets the minimum character count (inclusive).
func WithMinimumLength(length int) OptionFunc {
	return func(s *String) {
		s.minimumLength = length
	}
}

// WithMaximumLength sets the maximum character count (inclusive).
func WithMaximumLength(length int) OptionFunc {
	return func(s *String) {
		s.maximumLength = length
	}
}
