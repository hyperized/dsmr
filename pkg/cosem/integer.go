package cosem

import (
	"errors"
	"fmt"
	"strconv"
)

// IntegerOption configures an Integer.
type IntegerOption func(i *Integer)

// Integer is a parsed COSEM integer value.
type Integer struct {
	length int
	value  int64
}

// NewInteger parses value as a decimal integer. If WithIntegerLength is set,
// the string must have exactly that many characters.
func NewInteger(value string, options ...IntegerOption) (*Integer, error) {
	i := &Integer{}
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

// WithIntegerLength constrains the input to exactly length decimal digits.
// A length of 0 (the default) disables the check.
func WithIntegerLength(length int) IntegerOption {
	return func(i *Integer) {
		i.length = length
	}
}

// Value returns the parsed int64 value.
func (i *Integer) Value() int64 {
	return i.value
}
