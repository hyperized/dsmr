package data

import (
	"errors"
	"fmt"
	"github.com/hyperized/dsmr/cosem"
	"github.com/hyperized/dsmr/obis"
	"regexp"
	"strings"
)

type Object struct {
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
	return fmt.Sprintf("Object [%#v]", object)
}

const splitLineExpression string = "([0-9]-[0-9]:[0-9]+\\.[0-9]+\\.[0-9]+)(\\(.*\\))"
const splitLineError string = "could not parse line"

var (
	splitLine = regexp.MustCompile(splitLineExpression)
)

func NewFromLine(line string) (Object, error) {
	match := splitLine.FindStringSubmatch(strings.TrimSpace(line))

	if len(match) != 3 {
		return Object{}, errors.New(splitLineError)
	}

	reference := obis.References[match[1]]

	return Object{
		ObisIdentifier: reference.Identifier,
		RawValue:       match[2],
		Description:    reference.Description,
		Data:           reference.Format,
		Unit:           reference.Unit,
	}, nil
}
