package {{.CommandLine.Command.Name}}

{{if .CommandLine.Command.API}}
import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

func Validate(in interface{}) (msg string, err error) {
	t := reflect.TypeOf(in).Elem()
	v := reflect.ValueOf(in).Elem()

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(TagNameValidate)

		if tag == "" || tag == "-" || tag == "_" || tag == " " {
			continue
		}

		params := strings.Split(tag, ",")
		for _, param := range params {
			fmt.Printf("Validating: %s - %s\n", v.Type().Field(i).Name, param)

			if param == ValidateStrRequired {
				if v.Field(i).Pointer() == 0 {
					return ValidateStrRequiredErr, nil
				}
			}

			switch v.Field(i).Elem().Type() {
			case TypeOfString:
				if vMsg := ValidateString(param, v.Field(i).Elem().String()); vMsg != "" {
					return vMsg, nil
				}
			case TypeOfInt:
				if vMsg := ValidateInt(param, int(v.Field(i).Elem().Int())); vMsg != "" {
					return vMsg, nil
				}
			case TypeOfFloat32:
				if vMsg := ValidateFloat32(param, float32(v.Field(i).Elem().Float())); vMsg != "" {
					return vMsg, nil
				}
			case TypeOfFloat64:
				if vMsg := ValidateFloat64(param, v.Field(i).Elem().Float()); vMsg != "" {
					return vMsg, nil
				}
			}
		}
	}

	return
}

func ValidateString(param, in string) (msg string) {
	k, v := getTagKV(param)

	switch k {
	case ValidateStrMaxLength:
		if len(in) > cast.ToInt(v) {
			return ValidateStrMaxLengthErr
		}
	case ValidateStrMinLength:
		if len(in) < cast.ToInt(v) {
			return ValidateStrMinLengthErr
		}
	case ValidateStrMustHaveChars:
		if !allCharsInStr(v, in) {
			return ValidateStrMustHaveCharsErr
		}
	case ValidateStrCantHaveChars:
		if strings.IndexAny(in, v) > -1 {
			return ValidateStrCantHaveCharsErr
		}
	case ValidateStrOnlyHaveChars:
		if !onlyCharsInStr(v, in) {
			return ValidateStrOnlyHaveCharsErr
		}
	}

	return
}

func ValidateInt(param string, in int) (msg string) {
	k, v := getTagKV(param)

	switch k {
	case ValidateStrGreaterThan:
		if in < cast.ToInt(v) {
			return ValidateStrGreaterThanErr
		}
	case ValidateStrLessThan:
		if in > cast.ToInt(v) {
			return ValidateStrLessThanErr
		}
	}

	return
}

func ValidateFloat32(param string, in float32) (msg string) {
	k, v := getTagKV(param)

	switch k {
	case ValidateStrGreaterThan:
		if in < cast.ToFloat32(v) {
			return ValidateStrGreaterThanErr
		}
	case ValidateStrLessThan:
		if in > cast.ToFloat32(v) {
			return ValidateStrLessThanErr
		}
	}

	return
}

func ValidateFloat64(param string, in float64) (msg string) {
	k, v := getTagKV(param)

	switch k {
	case ValidateStrGreaterThan:
		if in < cast.ToFloat64(v) {
			return ValidateStrGreaterThanErr
		}
	case ValidateStrLessThan:
		if in > cast.ToFloat64(v) {
			return ValidateStrLessThanErr
		}
	}

	return
}
{{end}}
