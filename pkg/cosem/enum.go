package cosem

import (
	"fmt"
	"strconv"
)

// Enum is a parsed COSEM enumeration value (uint8, 0–255).
type Enum struct {
	value uint8
}

// NewEnum parses value as a decimal uint8. Values outside [0,255] or
// non-numeric inputs return an error.
func NewEnum(value string) (*Enum, error) {
	parsed, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid enum value %q: %w", value, err)
	}
	return &Enum{value: uint8(parsed)}, nil
}

// Value returns the parsed uint8 enumeration value.
func (e *Enum) Value() uint8 { return e.value }
