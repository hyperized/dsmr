package telegram

import "strconv"

// CRC-16/IBM (also known as CRC-16/ARC): polynomial 0x8005, reflected input
// and output, initial value 0x0000, no output XOR — equivalent to iterating
// with the reflected polynomial 0xA001. The DSMR P1 companion standard
// requires this algorithm over all bytes from and including the opening '/'
// up to and including the '!'.

// crcTable holds the precomputed CRC-16/IBM lookup for each possible input
// byte XORed into the running CRC value.
var crcTable = func() [256]uint16 {
	var t [256]uint16
	for i := range 256 {
		c := uint16(i)
		for range 8 {
			if c&0x0001 != 0 {
				c = (c >> 1) ^ 0xA001
			} else {
				c >>= 1
			}
		}
		t[i] = c
	}
	return t
}()

// ComputeCRC calculates the CRC-16/IBM checksum over data.
func ComputeCRC(data []byte) uint16 {
	var crc uint16
	for _, b := range data {
		crc = (crc >> 8) ^ crcTable[byte(crc)^b]
	}
	return crc
}

// ValidCRC returns true when the CRC of data matches expected (4 hex chars).
func ValidCRC(data []byte, expected string) bool {
	if len(expected) != 4 {
		return false
	}
	parsed, err := strconv.ParseUint(expected, 16, 16)
	if err != nil {
		return false
	}
	return ComputeCRC(data) == uint16(parsed)
}
