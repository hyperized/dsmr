package main

import (
	"github.com/hyperized/dsmr/telegram"
	"go.bug.st/serial"

	"log"
)

func main() {
	// P1
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	// TODO: Automatic port detection based on VID
	port, err := serial.Open("/dev/ttyUSB0", mode)
	defer func(port serial.Port) {
		err := port.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(port)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", port)

	parser := telegram.New(port)

	log.Printf("%#v\n", parser)

	_ = parser.Parse()

	// TODO: signal handling
}
