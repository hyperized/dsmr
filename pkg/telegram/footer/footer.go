// Package footer models the DSMR P1 telegram footer line ("!XXXX") which
// carries the CRC-16 checksum.
package footer

// OptionsFunc configures a Footer.
type OptionsFunc func(f *Footer)

// Footer holds the CRC string from a DSMR telegram's "!XXXX" line.
type Footer struct {
	crc string
}

// New creates a Footer, applying any provided options.
func New(options ...OptionsFunc) *Footer {
	h := &Footer{}

	for _, o := range options {
		o(h)
	}

	return h
}

func (f *Footer) String() string {
	return "Footer [CRC: " + f.crc + "]"
}

// CRC returns the four-character hex CRC string from the footer line.
func (f *Footer) CRC() string {
	return f.crc
}

// WithCRC sets the CRC string on the footer.
func WithCRC(crc string) OptionsFunc {
	return func(f *Footer) {
		f.crc = crc
	}
}
