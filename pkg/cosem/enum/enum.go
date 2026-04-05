// Package enum implements the COSEM Enum type (tag 22).
// An enum value is an unsigned 8-bit integer (range 0–255) represented as a
// decimal string in the telegram.
package enum

import (
	"fmt"
	"strconv"

	"github.com/hyperized/dsmr/pkg/cosem"
)

// OptionFunc configures an Enum.
type OptionFunc func(e *Enum)

// Enum holds a parsed COSEM enumeration value (uint8) with its type tag.
type Enum struct {
	tag   cosem.Cosem
	value uint8
}

// New parses value as a decimal uint8 (0–255). Values outside that range or
// non-numeric inputs return an error.
func New(value string, options ...OptionFunc) (*Enum, error) {
	e := &Enum{
		tag: cosem.Enum,
	}

	for _, o := range options {
		o(e)
	}

	parsed, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid enum value %q: %w", value, err)
	}

	e.value = uint8(parsed)
	return e, nil
}

// WithTag overrides the default COSEM tag for this enum.
func WithTag(tag cosem.Cosem) OptionFunc {
	return func(e *Enum) {
		e.tag = tag
	}
}

// Value returns the parsed uint8 enumeration value.
func (e *Enum) Value() uint8 {
	return e.value
}
