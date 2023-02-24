package template

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func MakePathAbosoluteOrExitOnError(path string) string {
	absPath, err := filepath.Abs(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return absPath
}

func CreateFile(relativeTemplPath string, data any, newFilePath string) {
	path, _ := os.Executable()
	generatorPath := filepath.Dir(path)
	templ, err := template.ParseFiles(filepath.Join(generatorPath, relativeTemplPath))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println("error occurred while creating a file", err)
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
