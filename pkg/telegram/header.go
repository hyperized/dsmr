package telegram

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

// NewHeader creates a Header from the four parsed identification fields.
func NewHeader(manufacturer string, baudRateID byte, model, version string) *Header {
	return &Header{
		manufacturer: manufacturer,
		baudRateID:   baudRateID,
		model:        model,
		version:      version,
	}
}

func (h *Header) String() string {
	return "Header [Manufacturer: " + h.manufacturer + ", Model: " + h.model + ", Version: " + h.version + "]"
}

// Manufacturer returns the 3-character manufacturer code (e.g. "ISk").
func (h *Header) Manufacturer() string { return h.manufacturer }

// BaudRateID returns the single-byte baud-rate identifier (normally '5').
func (h *Header) BaudRateID() byte { return h.baudRateID }

// Model returns the device model string from the identification field.
func (h *Header) Model() string { return h.model }

// Version returns the optional version suffix from the identification field.
func (h *Header) Version() string { return h.version }
