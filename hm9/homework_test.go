package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(p Person) string {
	var (
		pType  = reflect.TypeOf(p)
		pValue = reflect.ValueOf(p)
		b      = strings.Builder{}
	)

	n := pType.NumField()
	for i := 0; i < n; i++ {
		pField := pType.Field(i)
		nameFile, omitempty, ok := parse(pField.Tag)
		if !ok {
			continue
		}
		fieldValue := pValue.Field(i)
		if omitempty && fieldValue.IsZero() {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(nameFile)
		b.WriteByte('=')
		b.WriteString(fmt.Sprintf("%v", fieldValue))
	}

	return b.String()
}

func parse(t reflect.StructTag) (field string, om bool, ok bool) {
	meta, ok := t.Lookup("properties")
	if !ok {
		return
	}

	parts := strings.Split(meta, ",")

	field = strings.TrimSpace(parts[0])
	if len(field) == 0 {
		ok = false
	}

	if len(parts) == 1 {
		return
	}

	om = strings.TrimSpace(parts[1]) == "omitempty"
	return
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
