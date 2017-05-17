package {{.CommandLine.Command.Name}}

{{if .CommandLine.Command.API}}
import (
	"crypto/rand"
	"reflect"
	"strings"
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

func dereferenceStringArray(in []*string) (out []string) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceIntArray(in []*int) (out []int) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceInt32Array(in []*int32) (out []int32) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceInt64Array(in []*int64) (out []int64) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceFloat32Array(in []*float32) (out []float32) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceFloat64Array(in []*float64) (out []float64) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

func dereferenceBoolArray(in []*bool) (out []bool) {
	for _, inP := range in {
		out = append(out, *inP)
	}
	return
}

{{end}}
