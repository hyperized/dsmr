package telegram

// Footer holds the CRC string from a DSMR telegram's "!XXXX" line.
type Footer struct {
	crc string
}

// NewFooter creates a Footer carrying the supplied CRC string.
func NewFooter(crc string) *Footer {
	return &Footer{crc: crc}
}

func (f *Footer) String() string {
	return "Footer [CRC: " + f.crc + "]"
}

// CRC returns the four-character hex CRC string from the footer line.
func (f *Footer) CRC() string { return f.crc }
