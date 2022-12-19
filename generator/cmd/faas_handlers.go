package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Astenna/Nubes/generator/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/cmd"
)

var handlersCmd = &cobra.Command{
	Use:   "handlers",
	Short: "Generates handlers' definitions for AWS lambda deployment",
	Long:  `Generates handlers' definitions for AWS lambda deployment based on types and repositories indicated by the path`,

	Run: func(cmd *cobra.Command, args []string) {
		typesPath, _ := cmd.Flags().GetString("types")
		repositoriesPath, _ := cmd.Flags().GetString("repositories")
		generationDestination, _ := cmd.Flags().GetString("output")
		moduleName, _ := cmd.Flags().GetString("module")

		nobjectTypes := parser.GetNobjectsDefinedInPack(typesPath)
		functions := parser.PrepareStateChangingHandlers(MakePathAbosoluteOrExitOnError(typesPath), moduleName, nobjectTypes)
		repositoryFunctions := parser.PrepareRepositoriesHandlers(MakePathAbosoluteOrExitOnError(repositoriesPath), moduleName, nobjectTypes)

		GenerateStateChangingHandlers(generationDestination, functions)
		GenerateRepositoriesHandlers(generationDestination, repositoryFunctions)
	},
}

func init() {
	rootCmd.AddCommand(handlersCmd)

	var typesPath string
	var repositoriesPath string
	var handlersPath string
	var moduleName string

	handlersCmd.Flags().StringVarP(&typesPath, "types", "t", ".", "path to directory with types")
	handlersCmd.Flags().StringVarP(&repositoriesPath, "repositories", "r", ".", "path to directory with repositories")
	handlersCmd.Flags().StringVarP(&handlersPath, "output", "o", ".", "path where directory with handlers will be created")
	handlersCmd.Flags().StringVarP(&moduleName, "module", "m", ".", "module name of the source project")

	cmd.Execute()
}

func GenerateStateChangingHandlers(path string, functions []parser.HandlerFunc) {
	templ := ParseOrExitOnError("templates/handler_template.go.tmpl")
	generationDestPath := MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "state-changes"))
	os.MkdirAll(generationDestPath, 0777)

	for _, f := range functions {
		file, err := os.Create(filepath.Join(generationDestPath, f.HandlerName+".go"))
		if err != nil {
			fmt.Println(err)
		}
		err = templ.Execute(file, f)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
	}
}

func GenerateRepositoriesHandlers(path string, functions []parser.RepositoryCustomHandlerFunc) {
	templ := ParseOrExitOnError("templates/custom_repo_template.go.tmpl")
	repositoriesDirectoryPath := MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "repositories"))
	os.MkdirAll(repositoriesDirectoryPath, 0777)

	for _, f := range functions {
		file, err := os.Create(filepath.Join(repositoriesDirectoryPath, f.OperationName+f.TypeName+".go"))
		if err != nil {
			fmt.Println(err)
		}
		err = templ.Execute(file, f)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
	}
}

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
