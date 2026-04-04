// Package crc provides CRC-16 computation for DSMR P1 telegram validation.
//
// The algorithm used is CRC-16/IBM (also known as CRC-16/ARC):
// polynomial 0x8005, reflected input and output, initial value 0x0000, no
// output XOR. This is equivalent to iterating with the reflected polynomial
// 0xA001. The DSMR P1 companion standard requires this algorithm over all
// bytes from and including the opening '/' up to and including the '!'.
package crc

import (
	"fmt"
	"strings"
)

// Compute calculates the CRC-16/IBM checksum over data.
func Compute(data []byte) uint16 {
	var crc uint16
	for _, b := range data {
		crc ^= uint16(b)
		for range 8 {
			if crc&0x0001 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

// Valid returns true when the CRC of data matches expected (4 uppercase hex chars).
func Valid(data []byte, expected string) bool {
	if len(expected) != 4 {
		return false
	}
	return fmt.Sprintf("%04X", Compute(data)) == strings.ToUpper(expected)
}
