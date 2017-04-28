
import (
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"reflect"
	"strings"

	"encoding/hex"

	"crypto/aes"

	"github.com/spf13/cast"
)

const (
	TagNameValidate             = "validate"
	TagNameTransform            = "transform"
	TransformStrEncrypt         = "encrypt"
	TransformStrDecrypt         = "decrypt"
	TransformStrHash            = "hash"
	TransformStrTruncate        = "truncate"
	TransformStrTrimChars       = "trimChars"
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
	dummyString  string
	dummyInt     int
	dummyFloat32 float32
	dummyFloat64 float64
	dummyBool    bool

	TypeOfString  = reflect.TypeOf(dummyString)
	TypeOfInt     = reflect.TypeOf(dummyInt)
	TypeOfFloat32 = reflect.TypeOf(dummyFloat32)
	TypeOfFloat64 = reflect.TypeOf(dummyFloat64)
	TypeOfBool    = reflect.TypeOf(dummyBool)
)

func ErrorJSON(msg string) (out string) {
	return `{"error":"` + msg + `"}`
}

func Validate(in interface{}) (ok bool, msg string, err error) {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

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
					return false, ValidateStrRequiredErr, nil
				}
			}

			switch v.Field(i).Elem().Type() {
			case TypeOfString:
				if vMsg := ValidateString(param, v.Field(i).Elem().String()); vMsg != "" {
					return false, vMsg, nil
				}
			case TypeOfInt:
				if vMsg := ValidateInt(param, int(v.Field(i).Elem().Int())); vMsg != "" {
					return false, vMsg, nil
				}
			case TypeOfFloat32:
				if vMsg := ValidateFloat32(param, float32(v.Field(i).Elem().Float())); vMsg != "" {
					return false, vMsg, nil
				}
			case TypeOfFloat64:
				if vMsg := ValidateFloat64(param, v.Field(i).Elem().Float()); vMsg != "" {
					return false, vMsg, nil
				}
			}
		}
	}

	return true, "", nil
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

func Transform(in interface{}) (out interface{}, err error) {
	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(TagNameTransform)

		if tag == "" || tag == "-" || tag == "_" || tag == " " {
			continue
		}

		params := strings.Split(tag, ",")
		for _, param := range params {
			fmt.Printf("Transforming: %s - %s\n", v.Type().Field(i).Name, param)

			switch v.Field(i).Elem().Type() {
			case TypeOfString:
				if err := TransformString(param, v.Field(i).Elem()); err != nil {
					return in, err
				}
			}
		}
	}

	return in, nil
}

func TransformString(param string, value reflect.Value) (err error) {
	k, _ := getTagKV(param)

	switch k {
	case TransformStrHash:
		hashBytes32 := sha256.Sum256([]byte(value.String()))
		value.SetString(hex.EncodeToString(hashBytes32[:]))
	case TransformStrEncrypt:
		if value.String() == "" {
			return
		}
		if err := EncryptReflectValue(value); err != nil {
			fmt.Println("Failed Encryption...")
			return err
		}
	case TransformStrDecrypt:
		if value.String() == "" {
			return
		}
		if err := DecryptReflectValue(value); err != nil {
			fmt.Println("Failed Decryption...")
			return err
		}
	}

	return
}

func EncryptReflectValue(value reflect.Value) (err error) {
	fmt.Println("DONT USE THIS KEY IN PRODUCTION.. FETCH KEY FROM PKI")
	key := []byte("AES256Key-32Characters1234567890")

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	nonce := []byte("DON'T USE ME")
	fmt.Println("DONT USE THIS NONCE IN PRODUCTION.. GENERATE AND STORE RANDOM ONE")
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	// nonce := make([]byte, 12)
	// if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	// 	return err
	// }

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	cipherBytes := aesgcm.Seal(nil, nonce, []byte(value.String()), nil)

	value.SetString(hex.EncodeToString(cipherBytes))
	return
}

func DecryptReflectValue(value reflect.Value) (err error) {
	fmt.Println("DONT USE THIS KEY IN PRODUCTION.. FETCH KEY FROM PKI")
	key := []byte("AES256Key-32Characters1234567890")
	ciphertext, err := hex.DecodeString(value.String())
	if err != nil {
		return err
	}

	nonce := []byte("DON'T USE ME")
	fmt.Println("DONT USE THIS NONCE IN PRODUCTION.. FETCH THE ONE FOR THIS ENTRY")

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	value.SetString(string(plaintext))
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
