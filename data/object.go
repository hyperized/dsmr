package data

import (
	"errors"
	"fmt"
	"github.com/hyperized/dsmr/cosem"
	"github.com/hyperized/dsmr/obis"
	"regexp"
	"strings"
)

const DateTimeFormat = "060102150405"

type Object struct {
	Name           string
	ObisIdentifier string
	RawValue       string
	Description    string
	Data           cosem.Format
	Unit           string
}

func (object Object) Value() string {
	return object.RawValue
}

func (object Object) String() string {
	return fmt.Sprintf(object.Name + ": " + object.RawValue)
}

const splitLineExpression string = "([0-9]-[0-9]:[0-9]+\\.[0-9]+\\.[0-9]+)(\\(.*\\))"
const splitLineError string = "could not parse line"
const identifierError string = "could not find matching OBIS reference"

var (
	splitLine = regexp.MustCompile(splitLineExpression)
)

func NewFromLine(line string) (Object, error) {
	var (
		object   Object
		splitErr = errors.New(splitLineError)
		matchErr = errors.New(identifierError)
		match    = splitLine.FindStringSubmatch(strings.TrimSpace(line))
	)

	if len(match) != 3 {
		return object, splitErr
	}

	reference := obis.References[match[1]]

	if reference.Identifier == "" {
		return object, matchErr
	}

	return Object{
		Name:           reference.Name,
		ObisIdentifier: reference.Identifier,
		RawValue:       match[2],
		Description:    reference.Description,
		Data:           reference.Format,
		Unit:           reference.Unit,
	}, nil
}
