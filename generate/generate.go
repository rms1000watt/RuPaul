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
	dirCmd       = "cmd"
	extTpl       = ".tpl"
)

var (
	errGenFail = errors.New("Failed generating project")
)

func Generate(cfg Config) {
	fmt.Println("Config:", cfg)

	if err := os.Mkdir(dirOut, os.ModePerm); err != nil {
		fmt.Printf("Directory: \"%s\" already exists. Continuing...\n", dirOut)
	}

	templates := []string{
		".gitignore.tpl",
		"main.go.tpl",
		"cmd.root.go.tpl",
		"Readme.md.tpl",
	}

	for _, template := range templates {
		if err := genFile(cfg, template); err != nil {
			fmt.Println("Failed generating template:", template, ":", err)
		}
	}
}

func genFile(cfg Config, templateName string) (err error) {
	templateNameArr := strings.Split(templateName, ".")
	if len(templateNameArr) < 3 {
		fmt.Println("Bad templateName provided:", templateName)
		return errGenFail
	}

	fileName := strings.Join(templateNameArr[len(templateNameArr)-3:len(templateNameArr)-1], ".")
	dirs := templateNameArr[:len(templateNameArr)-3]
	dirs = append([]string{dirOut}, dirs...)

	templatePath := filepath.Join(dirTemplates, templateName)
	fileBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		fmt.Println("Failed reading:", err)
		return errGenFail
	}

	t, err := template.New(templateName).Parse(string(fileBytes))
	if err != nil {
		fmt.Println("Failed parsing template:", err)
		return errGenFail
	}

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

	return nil
}
