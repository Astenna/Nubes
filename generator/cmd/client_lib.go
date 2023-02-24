package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Astenna/Nubes/generator/parser"
	tp "github.com/Astenna/Nubes/generator/template_parser"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/cmd"
	"golang.org/x/exp/maps"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generates client project",
	Long:  `Generates client project based on types and repositories.`,

	Run: func(cmd *cobra.Command, args []string) {
		typesPath, _ := cmd.Flags().GetString("types")
		output, _ := cmd.Flags().GetString("output")
		projectName, _ := cmd.Flags().GetString("project-name")

		typesParser, err := parser.NewClientTypesParser(tp.MakePathAbosoluteOrExitOnError(typesPath))
		if err != nil {
			fmt.Println("Fatal occurred initialising type spec parser: %w", err)
			os.Exit(1)
		}
		typesParser.Run()
		definedTypes := maps.Values(typesParser.DefinedTypes)

		outputDirectoryPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(output, projectName))
		os.MkdirAll(outputDirectoryPath, 0777)

		path, _ := os.Executable()
		generatorPath := filepath.Dir(path)
		lambdaClient := tp.ParseOrExitOnError(tp.MakePathAbosoluteOrExitOnError(filepath.Join(generatorPath, "templates/client_lib/lambda_client.go.tmpl")))
		lambdaClientTemplInput := struct{ PackageName string }{PackageName: projectName}
		tp.CreateFileFromTemplate(lambdaClient, lambdaClientTemplInput, filepath.Join(outputDirectoryPath, "lambda_client.go"))

		templ := tp.ParseOrExitOnError(filepath.Join(generatorPath, "templates/client_lib/type.go.tmpl"))
		for _, typeDefinition := range definedTypes {
			typeDefinition.PackageName = projectName
			filePath := filepath.Join(outputDirectoryPath, typeDefinition.TypeNameLower+".go")
			tp.CreateFileFromTemplate(templ, typeDefinition, filePath)
			tp.RunGoimportsOnFile(filePath)
		}

		stub_templ := tp.ParseOrExitOnError(filepath.Join(generatorPath, "templates/client_lib/type_stubes.go.tmpl"))
		filePath := filepath.Join(outputDirectoryPath, "stubs.go")
		tp.CreateFileFromTemplate(stub_templ, struct {
			PackageName string
			Types       []*parser.StructTypeDefinition
		}{PackageName: projectName, Types: definedTypes}, filePath)
		tp.RunGoimportsOnFile(filePath)

		custom_ctors_templ := tp.ParseOrExitOnError(filepath.Join(generatorPath, "templates/client_lib/custom_ctors.go.tmpl"))
		customCtorTemplInput := struct {
			PackageName string
			CustomCtors []parser.CustomCtorDefinition
		}{PackageName: projectName, CustomCtors: typesParser.CustomCtorDefinitions}
		filePath = filepath.Join(outputDirectoryPath, "custom_ctors.go")
		tp.CreateFileFromTemplate(custom_ctors_templ, customCtorTemplInput, filePath)
		tp.RunGoimportsOnFile(filePath)

		other_decls_templ := tp.ParseOrExitOnError(filepath.Join(generatorPath, "templates/client_lib/other_decls.go.tmpl"))
		othetDeclsTemplInput := struct {
			PackageName string
			OtherDecls  parser.OtherDecls
		}{PackageName: projectName, OtherDecls: typesParser.OtherDecls}
		filePath = filepath.Join(outputDirectoryPath, "other_decls.go")
		tp.CreateFileFromTemplate(other_decls_templ, othetDeclsTemplInput, filePath)

		reference_templ := tp.ParseOrExitOnError(filepath.Join(generatorPath, "templates/client_lib/reference.go.tmpl"))
		referenceTmplInput := struct {
			PackageName string
		}{PackageName: projectName}
		filePath = filepath.Join(outputDirectoryPath, "reference.go")
		tp.CreateFileFromTemplate(reference_templ, referenceTmplInput, filePath)

		reference_nav_templ := tp.ParseOrExitOnError(filepath.Join(generatorPath, "templates/client_lib/reference_navigation_list.go.tmpl"))
		filePath = filepath.Join(outputDirectoryPath, "reference_navigation_list.go")
		tp.CreateFileFromTemplate(reference_nav_templ, referenceTmplInput, filePath)
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
