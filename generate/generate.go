package generate

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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
)

type HelperTemplate struct {
	Name     string
	Template string
}

func Generate(cfg Config) {
	fmt.Println("Config:", cfg)

	if err := os.Mkdir(dirOut, os.ModePerm); err != nil {
		fmt.Printf("Directory: \"%s\" already exists. Continuing...\n", dirOut)
	}

	templateArrs := [][]string{
		[]string{".gitignore.tpl"},
		[]string{"main.go.tpl"},
		[]string{"cmd.root.go.tpl", filepath.Join(dirHelpers, "cmd.root.globalArgs.tpl")},
		[]string{"Readme.md.tpl"},
	}

	for _, templateArr := range templateArrs {
		if err := genFile(cfg, templateArr); err != nil {
			fmt.Println("Failed generating templateArr:", templateArr, ":", err)
		}
	}
}

// func genFile(cfg Config, templateName string) (err error) {
func genFile(cfg Config, templateArr []string) (err error) {
	if len(templateArr) == 0 {
		fmt.Println("No templates in template arr")
		return errGenFail
	}

	templateName := templateArr[0]

	templateNameArr := strings.Split(templateName, ".")
	if len(templateNameArr) < 3 {
		fmt.Println("Bad templateName provided:", templateName)
		return errGenFail
	}

	fileName := strings.Join(templateNameArr[len(templateNameArr)-3:len(templateNameArr)-1], ".")
	dirs := templateNameArr[:len(templateNameArr)-3]
	dirs = append([]string{dirOut}, dirs...)

	helperTemplates := []HelperTemplate{}
	if len(templateArr) > 1 {
		for i := 1; i < len(templateArr); i++ {
			helperTplPath := filepath.Join(dirTemplates, templateArr[i])
			helperTplName := filepath.Base(templateArr[i])

			helperBytes, err := ioutil.ReadFile(helperTplPath)
			if err != nil {
				fmt.Println("Failed reading helper file:", err)
				return errGenFail
			}

			h, err := template.New(helperTplName).Parse(string(helperBytes))
			if err != nil {
				fmt.Println("Failed parsing helper template:", err)
				return errGenFail
			}

			var helperBytesBuffer bytes.Buffer
			if err := h.Execute(&helperBytesBuffer, cfg); err != nil {
				fmt.Println("Failed executing helper template:", err)
				return errGenFail
			}

			helperTemplates = append(helperTemplates, HelperTemplate{
				Name:     helperTplName,
				Template: string(helperBytesBuffer.Bytes()),
			})
		}
	}

	templatePath := filepath.Join(dirTemplates, templateName)
	t, err := template.New(templateName).ParseFiles(templatePath)
	if err != nil {
		fmt.Println("Failed parsing template:", err)
		return errGenFail
	}

	for _, helperTemplate := range helperTemplates {
		if _, err = t.New(helperTemplate.Name).Parse(helperTemplate.Template); err != nil {
			fmt.Println("Failed parsing helper template 2:", err)
			return errGenFail
		}
	}

	// TODO: Add helper functions.. ToCamelCase, FormatStringArg (checks type and adds double quotes)

	var outBuffer bytes.Buffer
	if err := t.Execute(&outBuffer, cfg); err != nil {
		fmt.Println("Failed executing template:", err)
		return errGenFail
	}

	if len(dirs) != 0 {
		if err := os.MkdirAll(filepath.Join(dirs...), os.ModePerm); err != nil {
			fmt.Println("Failed mkdir on dirs:", dirs, ":", err)
			return errGenFail
		}
	}

	dirsStr := filepath.Join(dirs...)
	completeFilePath := filepath.Join(dirsStr, fileName)

	if _, err := os.Stat(completeFilePath); err == nil {
		fmt.Println("NO overwrite, file exists:", completeFilePath)
		return nil
	}

	fmt.Println("Writing:", completeFilePath)
	if err := ioutil.WriteFile(completeFilePath, outBuffer.Bytes(), os.ModePerm); err != nil {
		fmt.Println("Failed writing file:", err)
		return errGenFail
	}

	// TODO: Run gofmt

	return nil
}
