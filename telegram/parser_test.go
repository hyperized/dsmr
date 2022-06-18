package telegram

import (
	"os"
	"testing"
)

func TestNewParserFromReader(t *testing.T) {
	file, err := os.Open("../examples/telegram_v5_0_2.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(file)

	parser := New(file)
	_ = parser.Parse()
	//t.Log(telegrams)
}

func TestHeaderFromToken(t *testing.T) {
	token, _ := tokenize("/ISk5\\2MT382-1000\n")
	header, ok := headerFromToken(token)
	if !ok {
		t.Error("no header found in string")
	}
	t.Log(header)
}

func TestDataFromToken(t *testing.T) {
	token, _ := tokenize("1-0:1.8.1(123456.789*kWh)\n")
	data, ok := dataFromToken(token)
	if !ok {
		t.Error("no data found in string")
	}
	t.Logf("%#v\n", data)
}

func TestFooterFromToken(t *testing.T) {
	token, _ := tokenize("!EF2F")
	header, ok := footerFromToken(token)
	if !ok {
		t.Error("no footer found in token")
	}
	t.Log(header)
}
