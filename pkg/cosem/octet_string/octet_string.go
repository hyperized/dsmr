// Package octet_string implements the COSEM OctetString type.
// Values are transmitted in telegrams as uppercase hex strings (2 hex chars per
// byte) and decoded to a raw []byte slice. Optional min/max length constraints
// are enforced on the decoded byte count.
package octet_string

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperized/dsmr/pkg/cosem"
)

// OptionFunc configures an OctetString.
type OptionFunc func(o *OctetString)

// OctetString holds a decoded COSEM octet string value with its type metadata.
type OctetString struct {
	tag           cosem.Cosem
	minimumLength int
	maximumLength int
	value         []byte
}

// New decodes value (a hex string) into bytes. An odd-length or non-hex input
// returns an error. Length constraints set via WithMinimumLength/WithMaximumLength
// are checked against the decoded byte count.
func New(value string, options ...OptionFunc) (*OctetString, error) {
	o := &OctetString{
		tag:           cosem.OctetString,
		minimumLength: 0,
		maximumLength: 0,
	}

	for _, opt := range options {
		opt(o)
	}

	if len(value)%2 != 0 {
		return nil, errors.New("hex string must have an even number of characters")
	}

	decoded, err := hex.DecodeString(strings.ToUpper(value))
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}

	octetLen := len(decoded)
	if octetLen < o.minimumLength || (o.maximumLength > 0 && octetLen > o.maximumLength) {
		return nil, fmt.Errorf("octet count out of range, expected %d to %d, found %d",
			o.minimumLength, o.maximumLength, octetLen)
	}

	o.value = decoded
	return o, nil
}

// WithTag overrides the default COSEM tag for this octet string.
func WithTag(tag cosem.Cosem) OptionFunc {
	return func(o *OctetString) {
		o.tag = tag
	}
}

// WithMinimumLength sets the minimum decoded byte count (inclusive).
func WithMinimumLength(length int) OptionFunc {
	return func(o *OctetString) {
		o.minimumLength = length
	}
}

// WithMaximumLength sets the maximum decoded byte count (inclusive; 0 = unlimited).
func WithMaximumLength(length int) OptionFunc {
	return func(o *OctetString) {
		o.maximumLength = length
	}
}

// Value returns the decoded byte slice.
func (o *OctetString) Value() []byte {
	return o.value
}

func (o *OctetString) String() string {
	return hex.EncodeToString(o.value)
}
