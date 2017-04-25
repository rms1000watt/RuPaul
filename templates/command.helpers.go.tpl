package {{.CommandLine.Command.Name}}

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	TagName = "validate"
	// TODO: Copy these same values from helpers.go
	ValidatorRequired = "required"
	TransformEncrypt  = "encrypt"
)

var (
	dummyString  string
	dummyInt     int
	dummyFloat32 float32
	dummyFloat64 float64

	TypeOfString  = reflect.TypeOf(dummyString)
	TypeOfInt     = reflect.TypeOf(dummyInt)
	TypeOfFloat32 = reflect.TypeOf(dummyFloat32)
	TypeOfFloat64 = reflect.TypeOf(dummyFloat64)
)

func ErrorJSON(msg string) (out string) {
	return `{"error":"` + msg + `"}`
}

func Validate(in interface{}) (ok bool, msg string, err error) {
	t := reflect.TypeOf(in)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(TagName)

		if tag == "" || tag == "-" || tag == "_" || tag == " " {
			continue
		}

		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
        fmt.Println(in)

		params := strings.Split(tag, ",")
		for _, param := range params {
			switch field.Type {
			case TypeOfString:
                // TODO: Fix this logic.. in is a struct not a string.
				inStr, interfaceOK := in.(string)
				if !interfaceOK {
					return false, "Failed handling interface (1)", nil
				}
				if vMsg := ValidateString(param, inStr); vMsg != "" {
					return false, vMsg, nil
				}
			case TypeOfInt:
				fmt.Println("Is Type Int!", param)
			}
		}
	}

	ok = true
	return
}

func ValidateString(param, in string) (msg string) {

	return
}
