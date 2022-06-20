package telegram

import (
	"bufio"
	"github.com/hyperized/dsmr/data"
	"io"
	"log"
	"strings"
)

type (
	Parser struct {
		scanner *bufio.Scanner
		tokens  []token
	}
	Telegram struct {
		header Header
		data   map[string]data.Object
		footer Footer
	}
	Footer struct {
		crc string
	}
	Header struct {
		manufacturer string
		model        string
		version      string
	}
)

func (telegram Telegram) String() string {
	var output []string

	output = append(output, "\nTelegram:\n"+telegram.header.String())
	for _, d := range telegram.data {
		output = append(output, d.String())
	}

	return strings.Join(output, "\n")
}

func (header *Header) String() string {
	return "Header [Manufacturer: " + header.manufacturer + ", Model: " + header.model + ", Version: " + header.version + "]"
}

func (footer Footer) String() string {
	return "Footer [CRC: " + footer.crc + "]"
}

func (parser *Parser) line() (Token, error, bool) {
	var (
		eof         = false
		tokenizeErr error
		scanErr     error
		token       Token
		ok          = parser.scanner.Scan()
	)
	if !ok {
		scanErr = parser.scanner.Err()
		if scanErr == nil {
			eof = true
		}
		return token, scanErr, eof
	}

	token, tokenizeErr = tokenize(parser.scanner.Text())

	return token, tokenizeErr, eof
}

func (parser *Parser) Parse() []Telegram {
	var (
		t  Telegram
		ts []Telegram
		c  = make(chan Telegram)
	)

	go parser.parseLines(c)

	for t = range c {
		ts = append(ts, t)
	}

	return ts
}

func (parser *Parser) ParseStream() {
	var (
		t Telegram
		c = make(chan Telegram)
	)

	go parser.parseLines(c)

	for t = range c {
		// TODO: Go routines to process provided telegrams
		go telegramToPrometheus(t)
	}
}

func (parser *Parser) parseLines(ch chan Telegram) {
	var (
		telegram    Telegram
		token       Token
		eof         bool
		tokenizeErr error
	)

	for {
		token, tokenizeErr, eof = parser.line()

		if tokenizeErr != nil {
			log.Println(tokenizeErr.Error())
			continue // Cannot parse the line
		}

		if eof {
			log.Println("found EOF, exiting loop")
			break
		}

		// Ignore empty lines
		if len(token.tokens) == 0 {
			continue
		}

		// Parse header
		if header, ok := headerFromToken(token); ok {
			telegram = Telegram{
				header: header,
				data:   make(map[string]data.Object),
			}
			continue
		}

		// OBIS Data
		if obis, ok := dataFromToken(token); ok {
			telegram.data[obis.Name] = obis
		}

		// Parse footer
		if footer, ok := footerFromToken(token); ok {
			// do CRC?
			telegram.footer = footer
			ch <- telegram
			continue
		}
	}

	log.Println("closing channel")
	close(ch)
}

func New(reader io.ReadWriteCloser) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(reader),
	}
}

func headerFromToken(token Token) (Header, bool) {
	var (
		literals []string
		h        Header
	)

	if token.tokens[0].kind != Slash {
		return h, false
	}

	// Grab literals
	for _, t := range token.tokens {
		if t.kind == Literal {
			literals = append(literals, t.value)
		}
	}

	// Return formed header
	return Header{
		manufacturer: literals[0],
		model:        literals[1],
		version:      literals[2],
	}, true
}

func dataFromToken(token Token) (data.Object, bool) {
	var (
		o   data.Object
		err error
	)

	o, err = data.NewFromLine(token.raw)
	if err != nil {
		return o, false
	}

	return o, true
}

func footerFromToken(token Token) (Footer, bool) {
	var (
		literal string
		f       Footer
	)

	if token.tokens[0].kind != Exclamation {
		return f, false
	}

	// Grab literal
	for _, t := range token.tokens {
		if t.kind == Literal {
			literal = t.value
		}
	}

	return Footer{
		crc: literal,
	}, true
}
