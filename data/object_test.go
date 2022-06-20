package data

import (
	"fmt"
	"testing"
)

func TestNewFromLine(t *testing.T) {
	var tests = []struct {
		input string
	}{
		{"1-0:32.32.0(00002)"},
		{"1-3:0.2.8(40)"},
		{"0-0:1.0.0(101209113020W)"},
		{"0-0:96.1.1(4B384547303034303436333935353037)"},
		{"1-0:1.8.1(123456.789*kWh)"},
		// No OBIS reference yet
		//{"1-0:99:97.0(2)(0:96.7.19)(101208152415W)(0000000240*s)(101208151004W)(00000000301*s)"},
		//{"0-0:96.13.0(303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F)"},
		//{"0-1:24.2.1(101209110000W)(12785.123*m3)"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s", tt.input)
		t.Run(name, func(t *testing.T) {
			object, err := NewFromLine(tt.input)
			t.Logf("%#v \n", object)
			if err != nil {
				t.Error(err)
			}
			if object.RawValue == "" {
				t.Error("could not parse value")
			}
			t.Log(object.Value[0])
		})
	}
}
