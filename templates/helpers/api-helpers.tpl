
import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

const (
	TagNameValidate             = "validate"
	TagNameTransform            = "transform"
	TagNameJSON                 = "json"
	TransformStrEncrypt         = "encrypt"
	TransformStrDecrypt         = "decrypt"
	TransformStrHash            = "hash"
	TransformStrPasswordHash    = "passwordHash"
	TransformStrTruncate        = "truncate"
	TransformStrTrimChars       = "trimChars"
	TransformStrTrimSpace       = "trimSpace"
	TransformStrDefault         = "default"
	ValidateStrMaxLength        = "maxLength"
	ValidateStrMinLength        = "minLength"
	ValidateStrGreaterThan      = "greaterThan"
	ValidateStrLessThan         = "lessThan"
	ValidateStrRequired         = "required"
	ValidateStrMustHaveChars    = "mustHaveChars"
	ValidateStrCantHaveChars    = "cantHaveChars"
	ValidateStrOnlyHaveChars    = "onlyHaveChars"
	ValidateStrMaxLengthErr     = "Failed Max Length Validation"
	ValidateStrMinLengthErr     = "Failed Min Length Validation"
	ValidateStrRequiredErr      = "Failed Required Validation"
	ValidateStrMustHaveCharsErr = "Failed Must Have Chars Validation"
	ValidateStrCantHaveCharsErr = "Failed Can't Have Chars Validation"
	ValidateStrOnlyHaveCharsErr = "Failed Only Have Chars Validation"
	ValidateStrGreaterThanErr   = "Failed Greater Than Validation"
	ValidateStrLessThanErr      = "Failed Less Than Validation"
)

var (
	dummyString   string
	dummyInt      int
	dummyInt64    int64
	dummyFloat32  float32
	dummyFloat64  float64
	dummyBool     bool
	dummyStringP  *string
	dummyIntP     *int
	dummyInt64P   *int64
	dummyFloat32P *float32
	dummyFloat64P *float64
	dummyBoolP    *bool

	TypeOfString   = reflect.TypeOf(dummyString)
	TypeOfInt      = reflect.TypeOf(dummyInt)
	TypeOfInt64    = reflect.TypeOf(dummyInt64)
	TypeOfFloat32  = reflect.TypeOf(dummyFloat32)
	TypeOfFloat64  = reflect.TypeOf(dummyFloat64)
	TypeOfBool     = reflect.TypeOf(dummyBool)
	TypeOfStringP  = reflect.TypeOf(dummyStringP)
	TypeOfIntP     = reflect.TypeOf(dummyIntP)
	TypeOfInt64P   = reflect.TypeOf(dummyInt64P)
	TypeOfFloat32P = reflect.TypeOf(dummyFloat32P)
	TypeOfFloat64P = reflect.TypeOf(dummyFloat64P)
	TypeOfBoolP    = reflect.TypeOf(dummyBoolP)
)

func ErrorJSON(msg string) (out string) {
	return `{"error":"` + msg + `"}`
}

func Unmarshal(r *http.Request, dst interface{}) (err error) {
	if r.Method == http.MethodGet {
		t := reflect.TypeOf(dst).Elem()
		v := reflect.ValueOf(dst).Elem()

		if err := r.ParseForm(); err != nil {
			return err
		}

		for i := 0; i < t.NumField(); i++ {
			jsonTag := t.Field(i).Tag.Get(TagNameJSON)
			jsonParams := strings.Split(jsonTag, ",")
			if len(jsonParams) == 0 {
				continue
			}
			jsonName := jsonParams[0]

			validateTag := t.Field(i).Tag.Get(TagNameValidate)
			validateParams := strings.Split(validateTag, ",")
			required := false
			for _, param := range validateParams {
				if param == ValidateStrRequired {
					required = true
				}
			}

			formValue := r.Form.Get(jsonName)
			if formValue == "" && required {
				return errors.New("Empty required field")
			}

			v.Field(i).Set(reflect.New(v.Field(i).Type().Elem()))

			switch v.Field(i).Type() {
			case TypeOfStringP:
				v.Field(i).Elem().SetString(formValue)
			case TypeOfIntP:
				fallthrough
			case TypeOfInt64P:
				v.Field(i).Elem().SetInt(cast.ToInt64(formValue))
			case TypeOfFloat64P:
				v.Field(i).Elem().SetFloat(cast.ToFloat64(formValue))
			case TypeOfFloat32P:
				fmt.Println("Float32 not supported")
				fallthrough
			default:
				fmt.Println("Field not set:", v.Type().Field(i).Name)
			}
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return err
	}

	return
}

func getRandomSalt() (salt []byte, err error) {
	salt = make([]byte, 32)
	_, err = rand.Read(salt)
	return
}

func getTagKV(param string) (k, v string) {
	paramArr := strings.Split(param, "=")

	k = paramArr[0]
	if len(paramArr) == 2 {
		v = paramArr[1]
	}
	return
}

func allCharsInStr(allChars, in string) (out bool) {
	for _, char := range allChars {
		if strings.Index(in, string(char)) == -1 {
			return
		}
	}
	return true
}

func onlyCharsInStr(onlyChars, in string) (out bool) {
	for _, char := range onlyChars {
		in = strings.Replace(in, string(char), "", -1)
	}
	return len(in) == 0
}
