package generate

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func yamlToTemplateCfg(cfg Config, commandName string) (sCfg TemplateConfig) {
	apiName := toUpperCamelCase(cfg.CommandLine.Commands[commandName].API)
	templateAPI := yamlToTemplateAPI(cfg.APIs[apiName], cfg)

	sCfg = TemplateConfig{
		Version:         cfg.Version,
		MainImportPath:  cfg.MainImportPath,
		CopyrightHolder: cfg.CopyrightHolder,
		API:             templateAPI,
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
		inputs := []Data{}
		for _, input := range yamlPath.Inputs {
			input = toUpperCamelCase(input)
			inputs = append(inputs, massageData(input, cfg.Datas[input]))
		}

		outputs := []Data{}
		for _, output := range yamlPath.Outputs {
			output = toUpperCamelCase(output)
			outputs = append(outputs, massageData(output, cfg.Datas[output]))
		}

		templatePaths = append(templatePaths, TemplatePath{
			Name:      yamlPath.Name,
			Pattern:   yamlPath.Pattern,
			Connector: yamlPath.Connector,
			Methods:   yamlPath.Methods,
			Inputs:    inputs,
			Outputs:   outputs,
		})
	}

	templateAPI = TemplateAPI{
		Name:          yamlAPI.Name,
		Type:          yamlAPI.Type,
		Serialization: yamlAPI.Serialization,
		Middlewares:   yamlAPI.Middlewares,
		Paths:         templatePaths,
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

	// Fill in Defaults
	if out.Type == "" {
		out.Type = typeString
	}

	return
}

// Courtesy of https://github.com/fatih/camelcase/blob/master/camelcase.go
func SplitCamel(src string) (entries []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
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
	return
}
