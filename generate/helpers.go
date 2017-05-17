package generate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cast"
)

const (
	TransformStrEncrypt      = "encrypt"
	TransformStrDecrypt      = "decrypt"
	TransformStrHash         = "hash"
	TransformStrPasswordHash = "passwordHash"
	TransformStrTruncate     = "truncate"
	TransformStrTrimChars    = "trimChars"
	TransformStrTrimSpace    = "trimSpace"
	TransformStrDefault      = "default"
	ValidateStrMaxLength     = "maxLength"
	ValidateStrMinLength     = "minLength"
	ValidateStrGreaterThan   = "greaterThan"
	ValidateStrLessThan      = "lessThan"
	ValidateStrRequired      = "required"
	ValidateStrMustHaveChars = "mustHaveChars"
	ValidateStrCantHaveChars = "cantHaveChars"
	ValidateStrOnlyHaveChars = "onlyHaveChars"
	DataTypeFloat32          = "float32"
	DataTypeFloat64          = "float64"
	DataTypeInt              = "int"
	DataTypeInt32            = "int32"
	DataTypeInt64            = "int64"
	DataTypeString           = "string"
	DataTypeBool             = "bool"
	DataTypeStringArr        = "[]string"
	DataTypeIntArr           = "[]int"
	DataTypeInt32Arr         = "[]int32"
	DataTypeInt64Arr         = "[]int64"
	DataTypeFloat32Arr       = "[]float32"
	DataTypeFloat64Arr       = "[]float64"
	DataTypeBoolArr          = "[]bool"
)

func yamlToTemplateCfg(cfg Config, commandName string) (sCfg TemplateConfig) {
	apiName := toUpperCamelCase(cfg.CommandLine.Commands[commandName].API)
	templateAPI := yamlToTemplateAPI(cfg.APIs[apiName], cfg)

	sCfg = TemplateConfig{
		Version:         cfg.Version,
		MainImportPath:  cfg.MainImportPath,
		CopyrightHolder: cfg.CopyrightHolder,
		API:             templateAPI,
		// Middlewares:     map[string]TemplateMiddleware{},
		CommandLine: TemplateCommandLine{
			AppName:             cfg.CommandLine.AppName,
			AppLongDescription:  cfg.CommandLine.AppLongDescription,
			AppShortDescription: cfg.CommandLine.AppShortDescription,
			GlobalArgs:          cfg.CommandLine.GlobalArgs,
			Command:             cfg.CommandLine.Commands[commandName],
		},
	}
	return
}

func yamlToTemplateAPI(yamlAPI API, cfg Config) (templateAPI TemplateAPI) {
	templatePaths := []TemplatePath{}
	for _, yamlPath := range yamlAPI.Paths {

		methods := []TemplateMethod{}
		for methodKey, methodValue := range yamlPath.Methods {
			inputs := []Data{}
			for _, input := range methodValue.Inputs {
				input = toUpperCamelCase(input)
				inputs = append(inputs, massageData(input, cfg.Datas[input]))
			}

			outputs := []Data{}
			for _, output := range methodValue.Outputs {
				output = toUpperCamelCase(output)
				outputs = append(outputs, massageData(output, cfg.Datas[output]))
			}

			methods = append(methods, TemplateMethod{
				Name:        strings.ToLower(methodKey),
				Inputs:      inputs,
				Outputs:     outputs,
				Middlewares: yamlToTemplateMiddlewares(methodValue.Middlewares, cfg),
				Connector:   methodValue.Connector,
			})
		}

		templatePaths = append(templatePaths, TemplatePath{
			Name:    yamlPath.Name,
			Pattern: yamlPath.Pattern,
			// Connector: yamlPath.Connector,
			Methods: methods,
			// Inputs:    inputs,
			// Outputs:   outputs,
		})
	}

	templateAPI = TemplateAPI{
		Name:            yamlAPI.Name,
		Type:            yamlAPI.Type,
		CertsPath:       yamlAPI.CertsPath,
		PubKeyFileName:  yamlAPI.PubKeyFileName,
		PrivKeyFileName: yamlAPI.PrivKeyFileName,
		Serialization:   yamlAPI.Serialization,
		Middlewares:     yamlToTemplateMiddlewares(yamlAPI.Middlewares, cfg),
		Paths:           templatePaths,
	}
	return
}

func yamlToTemplateMiddlewares(middlewares []string, cfg Config) (templateMiddlewares map[string]TemplateMiddleware) {
	templateMiddlewares = map[string]TemplateMiddleware{}

	for _, mw := range middlewares {
		templateMiddlewares[mw] = TemplateMiddleware{
			Options: getMiddlewareOptions(mw, cfg),
		}
	}
	return
}

func getMiddlewareOptions(mw string, cfg Config) (options []KV) {
	//TODO: There will probably be some optional options where defaults will be filled in here

	if len(cfg.Middlewares[mw].Options) != 0 {
		return cfg.Middlewares[mw].Options
	}
	return
}

func toUpperCamelCase(in string) (out string) {
	if in == "" {
		return
	}

	out = in

	separators := []string{
		" ",
		"_",
		"-",
	}
	for _, separator := range separators {
		nameArr := strings.Split(out, separator)
		for i := 0; i < len(nameArr); i++ {
			nameArr[i] = strings.Title(nameArr[i])
		}
		out = strings.Join(nameArr, "")
	}

	return
}

func massageData(name string, in Data) (out Data) {
	out = in
	out.Name = name

	if in.DisplayName == "" {
		out.DisplayName = ToSnakeCase(name)
	} else {
		out.DisplayName = in.DisplayName
	}

	// Fill in Defaults
	if out.Type == "" {
		out.Type = typeString
	}

	return
}

func GenTransformStr(in Data) (out string) {
	// TODO: Put all of these into constants and import in helpers.tpl
	if in.Encrypt {
		out += TransformStrEncrypt + ","
	}

	if in.Decrypt {
		out += TransformStrDecrypt + ","
	}

	if in.Hash {
		out += TransformStrHash + ","
	}

	if in.PasswordHash {
		out += TransformStrPasswordHash + ","
	}

	if in.TrimSpace {
		out += TransformStrTrimSpace + ","
	}

	if in.Truncate > 0 {
		out += TransformStrTruncate + "=" + cast.ToString(in.Truncate) + ","
	}

	if in.TrimChars != "" {
		out += TransformStrTrimChars + "=" + in.TrimChars + ","
	}

	if in.Default != "" {
		out += TransformStrDefault + "=" + in.Default + ","
	}

	return strings.Trim(out, ",")
}

func GenValidationStr(in Data) (out string) {
	// TODO: Put all of these into constants and import in helpers.tpl
	if in.Required {
		out += ValidateStrRequired + ","
	}

	if in.MaxLength > 0 {
		out += ValidateStrMaxLength + "=" + cast.ToString(in.MaxLength) + ","
	}

	if in.MinLength > 0 {
		out += ValidateStrMinLength + "=" + cast.ToString(in.MinLength) + ","
	}

	if in.MustHaveChars != "" {
		out += ValidateStrMustHaveChars + "=" + in.MustHaveChars + ","
	}

	if in.CantHaveChars != "" {
		out += ValidateStrCantHaveChars + "=" + in.CantHaveChars + ","
	}

	if in.OnlyHaveChars != "" {
		out += ValidateStrOnlyHaveChars + "=" + in.OnlyHaveChars + ","
	}

	if (in.Type == typeFloat || in.Type == typeInt) && in.GreaterThan != nil {
		out += ValidateStrGreaterThan + "=" + cast.ToString(in.GreaterThan) + ","
	}

	if (in.Type == typeFloat || in.Type == typeInt) && in.LessThan != nil {
		out += ValidateStrLessThan + "=" + cast.ToString(in.LessThan) + ","
	}

	return strings.Trim(out, ",")
}

func HandleQuotes(value, typeStr string) string {
	if strings.ToLower(typeStr) == typeString {
		return `"` + value + `"`
	}
	return value
}

func NormalizeConfig() {

}

func OutputInInputs(outputName string, inputs []Data) bool {
	for _, input := range inputs {
		if input.Name == outputName {
			return true
		}
	}
	return false
}

func EmptyValue(dataType string) (out string) {
	if dataType[:2] == "[]" {
		return dataType + "{}"
	}

	switch dataType {
	case DataTypeString:
		return "\"\""
	case DataTypeInt:
		fallthrough
	case DataTypeInt32:
		fallthrough
	case DataTypeInt64:
		return "0"
	case DataTypeFloat32:
		fallthrough
	case DataTypeFloat64:
		return "0.0"
	case DataTypeBool:
		return "false"
	}
	fmt.Println("DATA TYPE NOT DEFINED:", dataType)
	return "\"\""
}

// Courtesy of https://github.com/etgryphon/stringUp/blob/master/stringUp.go
var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func ToCamelCase(src string) (out string) {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	out = string(bytes.Join(chunks, nil))
	out = strings.ToLower(string(out[0])) + string(out[1:])
	return out
}

// Courtesy of https://github.com/fatih/camelcase/blob/master/camelcase.go
func ToSnakeCase(src string) (out string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	entries := []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}

	return strings.ToLower(strings.Join(entries, "_"))
}

func RemoveUnusedFile(completeFilePath string) {
	fileBytes, err := ioutil.ReadFile(completeFilePath)
	if err != nil {
		// Fail silently.. not a big deal
		return
	}

	// if !bytes.Contains(bytes.TrimSpace(fileBytes), []byte("\n")) && bytes.Equal(fileBytes[:7], []byte("package")) {
	if !bytes.Contains(bytes.TrimSpace(fileBytes), []byte("\n")) {
		fmt.Println("Removing:", completeFilePath)
		if err := os.Remove(completeFilePath); err != nil {
			// Fail silently.. not a big deal
			return
		}
	}
}

func GetHTTPMethod(method string) (httpMethod string) {
	method = strings.ToLower(method)
	switch method {
	case "connect":
		return "MethodConnect"
	case "delete":
		return "MethodDelete"
	case "get":
		return "MethodGet"
	case "head":
		return "MethodHead"
	case "options":
		return "MethodOptions"
	case "patch":
		return "MethodPatch"
	case "post":
		return "MethodPost"
	case "put":
		return "MethodPut"
	case "trace":
		return "MethodTrace"
	}

	fmt.Println("BAD METHOD PROVIDED!!", method)
	return
}

func CopyCertsPath(cfg Config) (certsPath string) {
	certsPaths := []string{}
	for _, api := range cfg.APIs {
		certsPaths = append(certsPaths, api.CertsPath)
	}

	if len(certsPaths) > 1 {
		prevCertsPath := certsPaths[0]
		for _, cp := range certsPaths {
			if prevCertsPath != cp {
				fmt.Println("Different Certs Path defined:", prevCertsPath, cp)
				return ""
			}
			prevCertsPath = cp
		}
		return dockerCopyCertsPath(certsPaths[0])
	}

	if len(certsPaths) == 1 {
		return dockerCopyCertsPath(certsPaths[0])
	}

	return ""
}

func dockerCopyCertsPath(certsPath string) string {
	if certsPath != "" {
		return "COPY " + certsPath + " /certs"
	}
	return ""
}

func FallbackSet(in, fallback string) string {
	if in != "" {
		return in
	}
	return fallback
}

func GetMethodMiddlewares(in string, cfg TemplateConfig) (out string) {
	mws := []string{}
	for _, path := range cfg.API.Paths {
		for _, method := range path.Methods {
			if method.Name == in {
				for k := range method.Middlewares {
					mws = append(mws, k)
				}
				break
			}
		}
	}

	out = ""
	for _, mw := range mws {
		out += ", Middleware" + mw
	}

	return
}

func GetPathMiddlewares(cfg TemplateConfig) (out string) {
	mws := []string{}
	for k := range cfg.API.Middlewares {
		mws = append(mws, k)
	}

	out = ""
	for _, mw := range mws {
		out += ", Middleware" + mw
	}
	return
}

func GetInputType(inputType string) (out string) {
	if len(inputType) < 2 {
		return inputType
	}
	if inputType[:2] != "[]" {
		return "*" + inputType
	}
	return "[]*" + inputType[2:]
}

func GetDereferenceFunc(outputType string) (out string) {
	switch outputType {
	case DataTypeStringArr:
		return "dereferenceStringArray"
	case DataTypeIntArr:
		return "dereferenceIntArray"
	case DataTypeInt32Arr:
		return "dereferenceInt32Array"
	case DataTypeInt64Arr:
		return "dereferenceInt64Array"
	case DataTypeFloat32Arr:
		return "dereferenceFloat32Array"
	case DataTypeFloat64Arr:
		return "dereferenceFloat64Array"
	case DataTypeBoolArr:
		return "dereferenceBoolArray"
	}

	return "*"
}

func GetProjectFolder(cfg Config) (projectFolder string) {
	splitPath := strings.Split(cfg.MainImportPath, "/")
	return splitPath[len(splitPath)-1]
}
