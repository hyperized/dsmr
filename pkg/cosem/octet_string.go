package cosem

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// OctetStringOption configures an OctetString.
type OctetStringOption func(o *OctetString)

// OctetString holds a decoded COSEM octet-string value (raw bytes from hex).
type OctetString struct {
	minimumLength int
	maximumLength int
	value         []byte
}

// NewOctetString decodes value (a hex string) into bytes. An odd-length or
// non-hex input returns an error. Length constraints set via
// WithMinimum/MaximumOctetStringLength are checked against the decoded byte count.
func NewOctetString(value string, options ...OctetStringOption) (*OctetString, error) {
	o := &OctetString{}
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

// WithMinimumOctetStringLength sets the minimum decoded byte count (inclusive).
func WithMinimumOctetStringLength(length int) OctetStringOption {
	return func(o *OctetString) { o.minimumLength = length }
}

// WithMaximumOctetStringLength sets the maximum decoded byte count (inclusive; 0 = unlimited).
func WithMaximumOctetStringLength(length int) OctetStringOption {
	return func(o *OctetString) { o.maximumLength = length }
}

// Value returns the decoded byte slice.
func (o *OctetString) Value() []byte { return o.value }

// String returns the hex-encoded representation of the decoded bytes.
func (o *OctetString) String() string { return hex.EncodeToString(o.value) }
