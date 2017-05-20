package generate

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/imdario/mergo"
	yaml "gopkg.in/yaml.v2"
)

const (
	dirOut       = "out"
	dirTemplates = "templates"
	dirHelpers   = "helpers"
	dirCmd       = "cmd"
	extTpl       = ".tpl"
)

var (
	errGenFail = errors.New("Failed generating project")
	funcMap    = template.FuncMap{
		"ToLower":              strings.ToLower,
		"Title":                strings.Title,
		"ToUpper":              strings.ToUpper,
		"TimeNowYear":          time.Now().Year,
		"GenValidationStr":     GenValidationStr,
		"GenTransformStr":      GenTransformStr,
		"HandleQuotes":         HandleQuotes,
		"ToSnakeCase":          ToSnakeCase,
		"ToCamelCase":          ToCamelCase,
		"OutputInInputs":       OutputInInputs,
		"NotOutputInInputs":    NotOutputInInputs,
		"EmptyValue":           EmptyValue,
		"GetHTTPMethod":        GetHTTPMethod,
		"CopyCertsPath":        CopyCertsPath,
		"FallbackSet":          FallbackSet,
		"GetMethodMiddlewares": GetMethodMiddlewares,
		"GetPathMiddlewares":   GetPathMiddlewares,
		"GetInputType":         GetInputType,
		"GetDereferenceFunc":   GetDereferenceFunc,
		"GetProjectFolder":     GetProjectFolder,
		"IsStruct":             IsStruct,
		"GetStructFields":      GetStructFields,
		"GetStructs":           GetStructs,
		"GetStructs2":          GetStructs2,
	}
)

type Template struct {
	TemplateName string
	FileName     string
	Data         interface{}
}

func Generate(cfg Config) {
	// TODO: Validate cfg so we don't have to ToLower everywhere in template

	fmt.Println("Config:", cfg)

	if err := os.Mkdir(dirOut, os.ModePerm); err != nil {
		fmt.Printf("Directory: \"%s\" already exists. Continuing...\n", dirOut)
	}

	helperFileNames, err := getHelperFileNames()
	if err != nil {
		fmt.Println("Failed reading helper files:", err)
		return
	}

	cfgs, err := importCfgs(cfg)
	if err != nil {
		fmt.Println("Failed importing other configs:", err)
		return
	}

	cfg, err = mergeConfigs(cfg, cfgs)
	if err != nil {
		fmt.Println("Failed merging configs:", err)
		return
	}

	templates := getTemplates(cfg)
	for _, tpl := range templates {
		if err := genFile(tpl, helperFileNames); err != nil {
			fmt.Println("Failed generating template:", tpl.TemplateName, ":", err)
		}
	}
}

func getHelperFileNames() (helperFileNames []string, err error) {
	helperFiles, err := ioutil.ReadDir(filepath.Join(dirTemplates, dirHelpers))
	if err != nil {
		return
	}

	for _, helperFile := range helperFiles {
		helperFileNames = append(helperFileNames, filepath.Join(dirHelpers, helperFile.Name()))
	}
	return
}

func importCfgs(cfg Config) (cfgs []Config, err error) {
	cfgs = []Config{}
	for _, filePath := range cfg.Imports {
		newCfg := &Config{}

		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Failed reading newCfg:", filePath)
			return nil, err
		}

		if err := yaml.Unmarshal(fileBytes, newCfg); err != nil {
			fmt.Println("Failed unmarshalling newCfg Yaml:", err)
			return nil, err
		}

		cfgs = append(cfgs, *newCfg)
	}
	return
}

func mergeConfigs(cfg Config, cfgs []Config) (Config, error) {
	for _, c := range cfgs {
		if err := mergo.Merge(&cfg, c); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

func getTemplates(cfg Config) (templates []Template) {
	singleTemplateNames := []string{
		".gitignore.tpl",
		"main.go.tpl",
		"cmd.root.go.tpl",
		"Readme.md.tpl",
		"License..tpl",
		"Dockerfile..tpl",
	}

	for _, singleTemplateName := range singleTemplateNames {
		templates = append(templates, Template{
			TemplateName: singleTemplateName,
			FileName:     singleTemplateName,
			Data:         cfg,
		})
	}

	for key, value := range cfg.CommandLine.Commands {
		lowerKey := strings.ToLower(key)
		templateCfg := yamlToTemplateCfg(cfg, key)

		templates = append(templates, Template{
			TemplateName: "cmd." + lowerKey + ".go.tpl",
			FileName:     "cmd.command.go.tpl",
			Data:         value,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + "." + lowerKey + ".go.tpl",
			FileName:     "command.command.go.tpl",
			Data:         templateCfg,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".config.go.tpl",
			FileName:     "command.config.go.tpl",
			Data:         templateCfg,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".middlewares.go.tpl",
			FileName:     "command.middlewares.go.tpl",
			Data:         templateCfg,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".helpers.go.tpl",
			FileName:     "command.helpers.go.tpl",
			Data:         templateCfg,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".unmarshal.go.tpl",
			FileName:     "command.unmarshal.go.tpl",
			Data:         templateCfg,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".validate.go.tpl",
			FileName:     "command.validate.go.tpl",
			Data:         templateCfg,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".transform.go.tpl",
			FileName:     "command.transform.go.tpl",
			Data:         templateCfg,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".data.go.tpl",
			FileName:     "command.data.go.tpl",
			Data:         templateCfg,
		})

		for _, apiPath := range templateCfg.API.Paths {
			ephemeralCfg := templateCfg
			ephemeralCfg.API.Paths = []TemplatePath{apiPath}

			templates = append(templates, Template{
				TemplateName: fmt.Sprintf("%s.%sHandler.go.tpl", lowerKey, apiPath.Name),
				FileName:     "command.handler.go.tpl",
				Data:         ephemeralCfg,
			})
		}
	}

	return
}

func genFile(tpl Template, helperFileNames []string) (err error) {
	templateName := tpl.TemplateName

	templateNameArr := strings.Split(templateName, ".")
	if len(templateNameArr) < 3 {
		fmt.Println("Bad templateName provided:", templateName)
		return errGenFail
	}

	templateFileName := filepath.Join(dirTemplates, tpl.FileName)
	fileBytes, err := ioutil.ReadFile(templateFileName)
	if err != nil {
		fmt.Println("Failed reading template file:", err)
		return
	}

	t, err := template.New(templateName).Funcs(funcMap).Parse(string(fileBytes))
	if err != nil {
		fmt.Println("Failed parsing template:", err)
		return errGenFail
	}

	fullHelperFileNames := []string{}
	for _, helperFileName := range helperFileNames {
		fullHelperFileNames = append(fullHelperFileNames, filepath.Join(dirTemplates, helperFileName))
	}

	if _, err := t.ParseFiles(fullHelperFileNames...); err != nil {
		fmt.Println("Failed parsing template:", err)
		return errGenFail
	}

	var outBuffer bytes.Buffer
	if err := t.Execute(&outBuffer, tpl.Data); err != nil {
		fmt.Println("Failed executing template:", err)
		return errGenFail
	}

	// Make the required directories for the project
	dirs := templateNameArr[:len(templateNameArr)-3]
	dirs = append([]string{dirOut}, dirs...)

	if len(dirs) != 0 {
		if err := os.MkdirAll(filepath.Join(dirs...), os.ModePerm); err != nil {
			fmt.Println("Failed mkdir on dirs:", dirs, ":", err)
			return errGenFail
		}
	}

	dirsStr := filepath.Join(dirs...)
	fileName := strings.Join(templateNameArr[len(templateNameArr)-3:len(templateNameArr)-1], ".")
	completeFilePath := filepath.Join(dirsStr, fileName)

	if filepath.Ext(completeFilePath) == "." {
		completeFilePath = completeFilePath[:len(completeFilePath)-1]
	}

	if _, err := os.Stat(completeFilePath); err == nil {
		fmt.Println("NO overwrite, file exists:", completeFilePath)
		return nil
	}

	fmt.Println("Writing:", completeFilePath)
	if err := ioutil.WriteFile(completeFilePath, outBuffer.Bytes(), os.ModePerm); err != nil {
		fmt.Println("Failed writing file:", err)
		return errGenFail
	}

	if filepath.Ext(completeFilePath) == ".go" {
		exec.Command("goimports", "-w", completeFilePath).CombinedOutput()
		exec.Command("gofmt", "-w", completeFilePath).CombinedOutput()
	}

	// TODO: Remove this--be smarter about which files to write
	RemoveUnusedFile(completeFilePath)

	return nil
}
