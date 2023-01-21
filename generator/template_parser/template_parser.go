package template_parser

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func ParseOrExitOnError(templatePath string) template.Template {
	templ, err := template.ParseFiles(templatePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return *templ
}

func MakePathAbosoluteOrExitOnError(path string) string {
	absPath, err := filepath.Abs(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return absPath
}

func CreateFileFromTemplate(templ template.Template, data any, newFilePath string) {
	file, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println("error occurred when creating a file", err)
	}
	err = templ.Execute(file, data)
	if err != nil {
		fmt.Println("template did not executed successfully", err)
	}
	file.Close()
}

type ExportTemplateInput struct {
	OrginalPackage            string
	OrginalPackageAlias       string
	IsNobjectInOrginalPackage map[string]bool
}
