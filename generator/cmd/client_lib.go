package cmd

import (
	"os"
	"path/filepath"

	"github.com/Astenna/Nubes/generator/parser"
	tp "github.com/Astenna/Nubes/generator/template_parser"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/cmd"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generates client project",
	Long:  `Generates client project based on types and repositories.`,

	Run: func(cmd *cobra.Command, args []string) {
		typesPath, _ := cmd.Flags().GetString("types")
		output, _ := cmd.Flags().GetString("output")
		projectName, _ := cmd.Flags().GetString("project-name")

		definedTypes, otherDecls := parser.ParsePackage(tp.MakePathAbosoluteOrExitOnError(typesPath))

		outputDirectoryPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(output, projectName))
		os.MkdirAll(outputDirectoryPath, 0777)

		lambdaClient := tp.ParseOrExitOnError("templates/client_lib/lambda_client.go.tmpl")
		lambdaClientInput := struct{ PackageName string }{PackageName: projectName}
		tp.CreateFileFromTemplate(lambdaClient, lambdaClientInput, filepath.Join(outputDirectoryPath, "lambda_client.go"))

		templ := tp.ParseOrExitOnError("templates/client_lib/type.go.tmpl")
		for _, typeDefinition := range definedTypes {
			typeDefinition.PackageName = projectName
			filePath := filepath.Join(outputDirectoryPath, typeDefinition.TypeNameLower+".go")
			tp.CreateFileFromTemplate(templ, typeDefinition, filePath)
			tp.RunGoimportsOnFile(filePath)
		}

		stub_templ := tp.ParseOrExitOnError("templates/client_lib/type_stubes.go.tmpl")
		filePath := filepath.Join(outputDirectoryPath, "stubs.go")
		tp.CreateFileFromTemplate(stub_templ, struct {
			PackageName string
			Types       []*parser.StructTypeDefinition
		}{PackageName: projectName, Types: definedTypes}, filePath)
		tp.RunGoimportsOnFile(filePath)

		repository_templ := tp.ParseOrExitOnError("templates/client_lib/repository.go.tmpl")
		tp.CreateFileFromTemplate(repository_templ, struct {
			PackageName string
		}{PackageName: projectName}, filepath.Join(outputDirectoryPath, "repository.go"))

		other_decls_templ := tp.ParseOrExitOnError("templates/client_lib/other_decls.go.tmpl")
		otherDecls.PackageName = projectName
		filePath = filepath.Join(outputDirectoryPath, "other_decls.go")
		tp.CreateFileFromTemplate(other_decls_templ, otherDecls, filePath)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	var typesPath string
	var outputPath string
	var projectName string

	clientCmd.Flags().StringVarP(&typesPath, "types", "t", ".", "path to directory with types")
	clientCmd.Flags().StringVarP(&outputPath, "output", "o", ".", "path where directory with client library will be created")
	clientCmd.Flags().StringVarP(&projectName, "project-name", "p", "client_lib", "name of the client library project")

	cmd.Execute()
}
