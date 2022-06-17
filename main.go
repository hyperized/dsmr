package main

import (
	"github.com/hyperized/dsmr/telegram"
	"log"
	"os"
)

func main() {
	file, err := os.Open("examples/telegram_v5_0_2.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	parser := telegram.NewFromReader(file)
	telegrams := parser.Parse()
	log.Println(telegrams)
}
