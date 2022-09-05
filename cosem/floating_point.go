package cosem

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"regexp"
	"strings"
)

type FloatingPoint struct {
	Tag         Cosem
	Length      int32
	MinDecimals int
	MaxDecimals int
	value       float64
}

func (f *FloatingPoint) WithString(value string) error {
	if !isValidString(value) {
		return errors.New("the provided string can only contain digits and dots")
	}

	digits, err := digitsAfterDot(value)
	if err != nil {
		return err
	}

	if digits < f.MinDecimals || digits > f.MaxDecimals {
		return fmt.Errorf("decimal point incorrect, expected to be between %d and %d, found %d",
			f.MinDecimals, f.MaxDecimals, digits,
		)
	}

	f.value = stringToFloat64(value, f.Length)
	return nil
}

func digitsAfterDot(value string) (int, error) {
	var splitMatchCount = 2
	split := strings.Split(value, ".")
	splitLen := len(split)

	// Since Split returns the original value if there's no split, overwrite splitLen
	if splitLen == 1 && split[0] == value {
		splitLen = 0
	}

	switch splitLen {
	case 0:
		return 0, nil
	case splitMatchCount:
		return len(split[1]), nil
	default:
		return 0, errors.New("multiple dots in string")
	}
}

func stringToFloat64(value string, decimals int32) float64 {
	dec, _ := decimal.NewFromString(value)
	dec.Round(decimals)
	result, _ := dec.Float64()
	return result
}

func isValidString(value string) bool {
	reg := regexp.MustCompile("[^0-9.]+")
	result := reg.FindString(value)
	return result == ""
}
