package crc_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hyperized/dsmr/pkg/telegram/crc"
	"github.com/stretchr/testify/assert"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected uint16
	}{
		{
			name:     "empty",
			input:    []byte{},
			expected: 0x0000,
		},
		{
			// CRC-16/IBM computed by hand for 0x31.
			name:     "single byte 0x31 ('1')",
			input:    []byte{0x31},
			expected: 0xD4C1,
		},
		{
			// Standard check value for CRC-16/IBM (CRC-16/ARC).
			name:     "standard check string '123456789'",
			input:    []byte("123456789"),
			expected: 0xBB3D,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := crc.Compute(tt.input)
			assert.Equal(t, tt.expected, got,
				"expected 0x%04X, got 0x%04X", tt.expected, got)
		})
	}
}

func TestValid(t *testing.T) {
	data := []byte("123456789")
	hex := fmt.Sprintf("%04X", crc.Compute(data))

	assert.True(t, crc.Valid(data, hex))
	assert.True(t, crc.Valid(data, "bb3d")) // case-insensitive
	assert.False(t, crc.Valid(data, "0000"))
	assert.False(t, crc.Valid(data, "BB3"))   // too short
	assert.False(t, crc.Valid(data, "BB3DE")) // too long
}

// TestExampleTelegramTwoCRC verifies the CRC-16/IBM checksum of the second
// example telegram with CRLF line endings.
//
// Note: the "EF2F" stored in the example file is a fictional placeholder used
// to exercise the parser's CRC-drop logic. The true CRC with CRLF endings is
// E47C.
func TestExampleTelegramTwoCRC(t *testing.T) {
	lines := []string{
		`/ISk5\2MT382-1000`,
		``,
		`1-3:0.2.8(50)`,
		`0-0:1.0.0(101209113020W)`,
		`0-0:96.1.1(4B384547303034303436333935353037)`,
		`1-0:1.8.1(123456.789*kWh)`,
		`1-0:1.8.2(123456.789*kWh)`,
		`1-0:2.8.1(123456.789*kWh)`,
		`1-0:2.8.2(123456.789*kWh)`,
		`0-0:96.14.0(0002)`,
		`1-0:1.7.0(01.193*kW)`,
		`1-0:2.7.0(00.000*kW)`,
		`0-0:96.7.21(00004)`,
		`0-0:96.7.9(00002)`,
		`1-0:99.97.0(2)(0-0:96.7.19)(101208152415W)(0000000240*s)(101208151004W)(0000000301*s)`,
		`1-0:32.32.0(00002)`,
		`1-0:52.32.0(00001)`,
		`1-0:72.32.0(00000)`,
		`1-0:32.36.0(00000)`,
		`1-0:52.36.0(00003)`,
		`1-0:72.36.0(00000)`,
		`0-0:96.13.0(303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F)`,
		`1-0:32.7.0(220.1*V)`,
		`1-0:52.7.0(220.2*V)`,
		`1-0:72.7.0(220.3*V)`,
		`1-0:31.7.0(001*A)`,
		`1-0:51.7.0(002*A)`,
		`1-0:71.7.0(003*A)`,
		`1-0:21.7.0(01.111*kW)`,
		`1-0:41.7.0(02.222*kW)`,
		`1-0:61.7.0(03.333*kW)`,
		`1-0:22.7.0(04.444*kW)`,
		`1-0:42.7.0(05.555*kW)`,
		`1-0:62.7.0(06.666*kW)`,
		`0-1:24.1.0(003)`,
		`0-1:96.1.0(3232323241424344313233343536373839)`,
		`0-1:24.2.1(101209112500W)(12785.123*m3)`,
	}

	var buf bytes.Buffer
	for _, l := range lines {
		buf.WriteString(l + "\r\n")
	}
	buf.WriteByte('!')

	assert.True(t, crc.Valid(buf.Bytes(), "E47C"),
		"CRC mismatch: computed %04X", crc.Compute(buf.Bytes()))
}

func TestRoundTrip(t *testing.T) {
	inputs := [][]byte{
		[]byte("/ISk5\\2MT382-1000\r\n\r\n1-3:0.2.8(50)\r\n!"),
		[]byte("hello world"),
		{0x00, 0xFF, 0xAA, 0x55},
	}
	for _, data := range inputs {
		crcVal := crc.Compute(data)
		assert.True(t, crc.Valid(data, fmt.Sprintf("%04X", crcVal)))
	}
}
