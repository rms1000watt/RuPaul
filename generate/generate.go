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
)

const (
	dirOut       = "out"
	dirTemplates = "templates"
	dirHelpers   = "helpers"
	dirCmd       = "cmd"
	extTpl       = ".tpl"
	typeString   = "string"
	typeBool     = "bool"
	typeFloat    = "float"
	typeInt      = "int"
)

var (
	errGenFail = errors.New("Failed generating project")
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
		templates = append(templates, Template{
			TemplateName: "cmd." + lowerKey + ".go.tpl",
			FileName:     "cmd.command.go.tpl",
			Data:         value,
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + "." + lowerKey + ".go.tpl",
			FileName:     "command.command.go.tpl",
			Data:         yamlToTemplateCfg(cfg, key),
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".config.go.tpl",
			FileName:     "command.config.go.tpl",
			Data:         yamlToTemplateCfg(cfg, key),
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".helpers.go.tpl",
			FileName:     "command.helpers.go.tpl",
			Data:         yamlToTemplateCfg(cfg, key),
		})
		templates = append(templates, Template{
			TemplateName: lowerKey + ".data.go.tpl",
			FileName:     "command.data.go.tpl",
			Data:         yamlToTemplateCfg(cfg, key),
		})
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

	funcMap := template.FuncMap{
		"ToLower":          strings.ToLower,
		"Title":            strings.Title,
		"TimeNowYear":      time.Now().Year,
		"GenValidationStr": GenValidationStr,
		"GenTransformStr":  GenTransformStr,
		"HandleQuotes":     HandleQuotes,
		"ToSnakeCase":      ToSnakeCase,
		"ToCamelCase":      ToCamelCase,
		"OutputInInputs":   OutputInInputs,
		"EmptyValue":       EmptyValue,
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

	RemoveUnusedFile(completeFilePath)

	return nil
}

func RemoveUnusedFile(completeFilePath string) {
	fileBytes, err := ioutil.ReadFile(completeFilePath)
	if err != nil {
		// Fail silently.. not a big deal
		return
	}

	if !bytes.Contains(bytes.TrimSpace(fileBytes), []byte("\n")) {
		if err := os.Remove(completeFilePath); err != nil {
			// Fail silently.. not a big deal
			return
		}
	}
}
