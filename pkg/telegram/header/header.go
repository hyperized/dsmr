// Package header models the identification line of a DSMR P1 telegram.
package header

// OptionsFunc configures a Header.
type OptionsFunc func(h *Header)

// Header represents the identification line of a DSMR P1 telegram, parsed
// according to IEC 62056-21.
//
// Format: /MFRb\identification
//   - MFR  – 3-character manufacturer identification
//   - b    – baud-rate identifier ('5' in all DSMR P1 implementations)
//   - identification – device model and serial number separated by '-'
type Header struct {
	manufacturer string
	baudRateID   byte
	model        string
	version      string
}

// New creates a Header, applying any provided options.
func New(options ...OptionsFunc) *Header {
	h := &Header{}

	for _, o := range options {
		o(h)
	}

	return h
}

func (h *Header) String() string {
	return "Header [Manufacturer: " + h.manufacturer + ", Model: " + h.model + ", Version: " + h.version + "]"
}

// Manufacturer returns the 3-character manufacturer code (e.g. "ISk").
func (h *Header) Manufacturer() string {
	return h.manufacturer
}

// BaudRateID returns the single-byte baud-rate identifier (normally '5').
func (h *Header) BaudRateID() byte {
	return h.baudRateID
}

// Model returns the device model string from the identification field.
func (h *Header) Model() string {
	return h.model
}

// Version returns the optional version suffix from the identification field.
func (h *Header) Version() string {
	return h.version
}

// WithManufacturer sets the manufacturer code on the header.
func WithManufacturer(manufacturer string) OptionsFunc {
	return func(h *Header) {
		h.manufacturer = manufacturer
	}
}

// WithBaudRateID sets the baud-rate identifier byte on the header.
func WithBaudRateID(id byte) OptionsFunc {
	return func(h *Header) {
		h.baudRateID = id
	}
}

// WithModel sets the device model string on the header.
func WithModel(model string) OptionsFunc {
	return func(h *Header) {
		h.model = model
	}
}

// WithVersion sets the version suffix on the header.
func WithVersion(version string) OptionsFunc {
	return func(h *Header) {
		h.version = version
	}
}
