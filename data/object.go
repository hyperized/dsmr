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
	Metric         obis.Metric
	ObisIdentifier string
	RawValue       string
	Value          []string
	Description    string
	Format         cosem.Format
	Unit           string
}

func (object Object) String() string {
	return fmt.Sprintf(object.Name + ": " + object.Value[0] + object.Unit)
}

const splitLineExpression string = "([0-9]-[0-9]:[0-9]+[\\.:][0-9]+\\.[0-9]+)\\((.*)\\)"
const splitLineError string = "object: could not parse line"
const identifierError string = "object: could not find matching OBIS reference"
const nullError string = "object: no value could be parsed from this line"

var (
	splitLine = regexp.MustCompile(splitLineExpression)
)

func NewFromLine(line string) (Object, error) {
	var (
		object   Object
		splitErr = errors.New(splitLineError)
		matchErr = errors.New(identifierError)
		nullErr  = errors.New(nullError)
		match    = splitLine.FindStringSubmatch(strings.TrimSpace(line))
	)

	if len(match) != 3 {
		return object, splitErr
	}

	reference := obis.References[match[1]]

	if reference.Identifier == "" {
		return object, matchErr
	}

	value := getValue(match[2])
	if value[0] == "" {
		return object, nullErr
	}

	return Object{
		Name:           reference.Name,
		Metric:         reference.Metric,
		ObisIdentifier: reference.Identifier,
		RawValue:       match[2],
		Value:          value,
		Description:    reference.Description,
		Format:         reference.Format,
		Unit:           reference.Unit,
	}, nil
}

func getValue(raw string) []string {
	var r []string
	values := strings.Split(raw, ")(")
	for _, v := range values {
		// Strip units out of it, we already know
		p := strings.Split(v, "*")
		r = append(r, p[0])

		// TODO: Parse timestamp
	}
	return r
}
