// Package integer implements the COSEM Integer COSEM type.
// It parses a decimal string into an int64 and optionally enforces a fixed
// digit-count constraint (e.g. "I4" requires exactly 4 digits).
package integer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperized/dsmr/pkg/cosem"
)

// OptionFunc configures an Integer.
type OptionFunc func(i *Integer)

// Integer holds a parsed COSEM integer value together with its type metadata.
type Integer struct {
	tag    cosem.Cosem
	length int
	value  int64
}

// New parses value as a decimal integer. If WithLength is set, the string must
// have exactly that many characters.
func New(value string, options ...OptionFunc) (*Integer, error) {
	i := &Integer{
		tag:    cosem.Integer,
		length: 0,
	}

	for _, o := range options {
		o(i)
	}

	if len(value) == 0 {
		return nil, errors.New("value must not be empty")
	}

	if i.length > 0 && len(value) != i.length {
		return nil, fmt.Errorf("integer length incorrect, expected %d digits, found %d", i.length, len(value))
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer value %q: %w", value, err)
	}

	i.value = parsed
	return i, nil
}

// WithTag overrides the default COSEM tag for this integer.
func WithTag(tag cosem.Cosem) OptionFunc {
	return func(i *Integer) {
		i.tag = tag
	}
}

// WithLength constrains the input to exactly length decimal digits.
// A length of 0 disables the check.
func WithLength(length int) OptionFunc {
	return func(i *Integer) {
		i.length = length
	}
}

// Value returns the parsed int64 value.
func (i *Integer) Value() int64 {
	return i.value
}
