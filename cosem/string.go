package cosem

import (
	"fmt"
)

type String struct {
	Tag       Cosem
	MinLength int
	MaxLength int
	value     string
}

func (s *String) WithString(value string) error {
	strLen := len(value)
	if strLen < s.MinLength || strLen > s.MaxLength {
		return fmt.Errorf("character count exceeds range, expected to be between %d and %d, found %d",
			s.MinLength, s.MaxLength, strLen,
		)
	}
	s.value = value
	return nil
}
