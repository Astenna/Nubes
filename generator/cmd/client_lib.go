package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Astenna/Nubes/generator/parser"
	templ "github.com/Astenna/Nubes/generator/template"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/cmd"
	"golang.org/x/exp/maps"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generates client project",
	Long:  `Generates client project based on types and repositories.`,

	Run: func(cmd *cobra.Command, _ []string) {
		typesPath, _ := cmd.Flags().GetString("types")
		output, _ := cmd.Flags().GetString("output")
		projectName, _ := cmd.Flags().GetString("project-name")

		typesParser, err := parser.NewClientTypesParser(templ.MakePathAbosoluteOrExitOnError(typesPath))
		if err != nil {
			fmt.Println("Fatal error occurred initialising type spec parser: %w", err)
			os.Exit(1)
		}
		typesParser.Run()

		outputDirectoryPath := templ.MakePathAbosoluteOrExitOnError(filepath.Join(output, projectName))
		os.MkdirAll(outputDirectoryPath, 0777)

		definedTypes := maps.Values(typesParser.DefinedTypes)
		for _, typeDefinition := range definedTypes {
			typeDefinition.PackageName = projectName
			filePath := filepath.Join(outputDirectoryPath, typeDefinition.TypeNameLower+".go")
			templ.CreateFile("template/client_lib/type.go.tmpl", typeDefinition, filePath)
			templ.RunGoimportsOnFile(filePath)
		}

		filePath := filepath.Join(outputDirectoryPath, "stubs.go")
		templ.CreateFile("template/client_lib/type_stubs.go.tmpl", struct {
			PackageName string
			Types       []*parser.StructTypeDefinition
		}{PackageName: projectName, Types: definedTypes}, filePath)
		templ.RunGoimportsOnFile(filePath)

		customCtorTemplInput := struct {
			PackageName string
			CustomCtors []parser.CustomCtorDefinition
		}{PackageName: projectName, CustomCtors: typesParser.CustomCtorDefinitions}
		filePath = filepath.Join(outputDirectoryPath, "custom_ctors.go")
		templ.CreateFile("template/client_lib/custom_ctors.go.tmpl", customCtorTemplInput, filePath)
		templ.RunGoimportsOnFile(filePath)

		othetDeclsTemplInput := struct {
			PackageName string
			OtherDecls  parser.OtherDecls
		}{PackageName: projectName, OtherDecls: typesParser.OtherDecls}
		filePath = filepath.Join(outputDirectoryPath, "other_decls.go")
		templ.CreateFile("template/client_lib/other_decls.go.tmpl", othetDeclsTemplInput, filePath)

		referenceTmplInput := struct {
			PackageName string
		}{PackageName: projectName}
		filePath = filepath.Join(outputDirectoryPath, "reference.go")
		templ.CreateFile("template/client_lib/reference.go.tmpl", referenceTmplInput, filePath)

		filePath = filepath.Join(outputDirectoryPath, "reference_navigation_list.go")
		templ.CreateFile("template/client_lib/reference_navigation_list.go.tmpl", referenceTmplInput, filePath)

		filePath = filepath.Join(outputDirectoryPath, "reference_ctors.go")
		templ.CreateFile("template/client_lib/reference_ctors.go.tmpl", struct {
			PackageName string
			Types       []*parser.StructTypeDefinition
		}{PackageName: projectName, Types: definedTypes}, filePath)

		lambdaClientTemplInput := struct{ PackageName string }{PackageName: projectName}
		templ.CreateFile("template/client_lib/lambda_client.go.tmpl", lambdaClientTemplInput, filepath.Join(outputDirectoryPath, "lambda_client.go"))
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	var typesPath string
	var outputPath string
	var projectName string

	clientCmd.Flags().StringVarP(&typesPath, "types", "t", ".", "path to package with types definitions")
	clientCmd.Flags().StringVarP(&outputPath, "output", "o", ".", "path where the directory with the client library will be created")
	clientCmd.Flags().StringVarP(&projectName, "project-name", "p", "client_lib", "name of the generated package")

	cmd.Execute()
}
