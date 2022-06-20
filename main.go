package main

import (
	"github.com/hyperized/dsmr/telegram"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.bug.st/serial"
	"log"
	"net/http"
)

func main() {
	// TODO CLI Flags & ENV variables
	adjustingSweepAngle()

	// Prometheus
	go func() {
		log.Println("starting prometheus webserver ...")
		http.Handle("/", http.RedirectHandler("/metrics", 302))
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("stopping prometheus webserver ...")
	}()

	// P1
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	// TODO: Automatic port detection based on VID
	log.Println("connecting to serial port ...")
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

	parser.ParseStream()

	// TODO: signal handling

	spacialControls()
}

func adjustingSweepAngle() {
	// https://youtu.be/gG0YJZvfUvo?t=365
	log.Println("")
	log.Println("Recalibrate azimuth sweep angle. Adjust elevation scan ..")
	log.Println("")
}

func spacialControls() {
	// https://youtu.be/sXMhGADyMxE?t=705
	log.Println("")
	log.Println("")
	log.Println("Well, that was some experience.")
	log.Println("Now .. just let me adjust the spacial controls and we'll move to another observation point.")
	log.Println("")
	log.Println("")
}
