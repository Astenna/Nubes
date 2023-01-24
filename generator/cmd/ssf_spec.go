package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Astenna/Nubes/generator/database"
	"github.com/Astenna/Nubes/generator/parser"
	tp "github.com/Astenna/Nubes/generator/template_parser"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/cmd"
)

var ssfSpecCmd = &cobra.Command{
	Use:   "handlers",
	Short: "Generates handlers' definitions for AWS lambda deployment",
	Long:  `Generates handlers' definitions for AWS lambda deployment based on types and repositories indicated by the path`,

	Run: func(cmd *cobra.Command, args []string) {
		typesPath, _ := cmd.Flags().GetString("types")
		generationDestination, _ := cmd.Flags().GetString("output")
		moduleName, _ := cmd.Flags().GetString("module")
		dbInit, _ := cmd.Flags().GetBool("dbInit")
		generateDeploymentFiles, _ := cmd.Flags().GetBool("deplFiles")

		typesPath = tp.MakePathAbosoluteOrExitOnError(typesPath)

		typeSpecParser, err := parser.NewTypeSpecParser(typesPath)
		if err != nil {
			fmt.Println("Fatal occurred initialising type spec parser: %w", err)
			os.Exit(1)
		}
		typeSpecParser.Run(moduleName)

		GenerateStateChangingHandlers(generationDestination, typeSpecParser.Handlers)
		GenerateGenericHandlers(generationDestination, typeSpecParser.Output)
		GenerateCustomConstructorsHandlers(generationDestination, typeSpecParser.CustomCtors)

		if generateDeploymentFiles {
			serviceName := lastString(strings.Split(moduleName, "/"))
			serverlessInput := ServerlessTemplateInput{ServiceName: serviceName, StateFuncs: typeSpecParser.Handlers, CustomCtors: typeSpecParser.CustomCtors}
			GenerateDeploymentFiles(generationDestination, serverlessInput)
		}

		if dbInit {
			database.CreateTypeTables(typeSpecParser.Output)
		}
	},
}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}

func init() {
	rootCmd.AddCommand(ssfSpecCmd)

	var typesPath string
	var handlersPath string
	var moduleName string
	var dbInit bool
	var generateDeploymentFiles bool

	ssfSpecCmd.Flags().StringVarP(&typesPath, "types", "t", ".", "path to directory with types")
	ssfSpecCmd.Flags().StringVarP(&handlersPath, "output", "o", ".", "path where directory with handlers will be created")
	ssfSpecCmd.Flags().StringVarP(&moduleName, "module", "m", "MISSING_MODULE_NAME", "module name of the source project")
	ssfSpecCmd.Flags().BoolVarP(&dbInit, "dbInit", "i", false, "boolean, indicates whether database should be initialized by creation of tables based on type names")
	ssfSpecCmd.Flags().BoolVarP(&generateDeploymentFiles, "deplFiles", "g", true, "boolean, indicates whether deployment files for AWS lambdas are to be created")

	cmd.Execute()
}

type ServerlessTemplateInput struct {
	ServiceName string
	StateFuncs  []parser.StateChangingHandler
	CustomCtors []parser.CustomCtorDefinition
}

func GenerateDeploymentFiles(path string, templateInput ServerlessTemplateInput) {
	serverlessTempl := tp.ParseOrExitOnError("templates/type_spec/deployment/serverless.yml.tmpl")
	fileName := filepath.Join(tp.MakePathAbosoluteOrExitOnError(path), "serverless.yml")
	tp.CreateFileFromTemplate(serverlessTempl, templateInput, fileName)

	buildScriptTempl := tp.ParseOrExitOnError("templates/type_spec/deployment/build_handlers.sh.tmpl")
	fileName = filepath.Join(tp.MakePathAbosoluteOrExitOnError(path), "build_handlers.sh")
	tp.CreateFileFromTemplate(buildScriptTempl, nil, fileName)
}

func GenerateStateChangingHandlers(path string, functions []parser.StateChangingHandler) {
	var handlerDir string
	var ownerHandlerNameCombined string
	templ := tp.ParseOrExitOnError("templates/type_spec/state_changing_template.go.tmpl")
	generationDestPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "state-changes"))

	for _, f := range functions {
		ownerHandlerNameCombined = f.ReceiverType + f.MethodName
		handlerDir = filepath.Join(generationDestPath, ownerHandlerNameCombined)
		os.MkdirAll(handlerDir, 0777)
		path = filepath.Join(handlerDir, ownerHandlerNameCombined+".go")
		tp.CreateFileFromTemplate(templ, f, path)
		tp.RunGoimportsOnFile(path)
	}
}

func GenerateGenericHandlers(path string, paredPkg parser.ParsedPackage) {
	templ := tp.ParseOrExitOnError("templates/type_spec/get_field_template.go.tmpl")
	generationDestPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "GetField"))
	os.MkdirAll(generationDestPath, 0777)
	getPath := filepath.Join(generationDestPath, "GetField.go")
	tp.CreateFileFromTemplate(templ, nil, getPath)

	templ = tp.ParseOrExitOnError("templates/type_spec/set_field_template.go.tmpl")
	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "SetField"))
	os.MkdirAll(generationDestPath, 0777)
	setPath := filepath.Join(generationDestPath, "SetField.go")
	tp.CreateFileFromTemplate(templ, nil, setPath)

	templ = tp.ParseOrExitOnError("templates/type_spec/load_template.go.tmpl")
	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "Load"))
	os.MkdirAll(generationDestPath, 0777)
	loadPath := filepath.Join(generationDestPath, "Load.go")
	tp.CreateFileFromTemplate(templ, nil, loadPath)

	templ = tp.ParseOrExitOnError("templates/type_spec/export_template.go.tmpl")
	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "Export"))
	os.MkdirAll(generationDestPath, 0777)
	exportPath := filepath.Join(generationDestPath, "Export.go")
	intput := tp.ExportTemplateInput{IsNobjectInOrginalPackage: paredPkg.IsNobjectInOrginalPackage, OrginalPackageAlias: parser.OrginalPackageAlias, OrginalPackage: paredPkg.ImportPath}
	tp.CreateFileFromTemplate(templ, intput, exportPath)

	templ = tp.ParseOrExitOnError("templates/type_spec/delete_template.go.tmpl")
	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "Delete"))
	os.MkdirAll(generationDestPath, 0777)
	deletePath := filepath.Join(generationDestPath, "Delete.go")
	tp.CreateFileFromTemplate(templ, nil, deletePath)

	templ = tp.ParseOrExitOnError("templates/type_spec/reference_get_by_index.go.tmpl")
	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "reference", "GetByIndex"))
	os.MkdirAll(generationDestPath, 0777)
	referenceIndexPath := filepath.Join(generationDestPath, "ReferenceGetByIndex.go")
	tp.CreateFileFromTemplate(templ, nil, referenceIndexPath)

	templ = tp.ParseOrExitOnError("templates/type_spec/reference_get_sort_key.go.tmpl")
	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "reference", "GetSortKeysByPartionKey"))
	os.MkdirAll(generationDestPath, 0777)
	referenceSortKeyPath := filepath.Join(generationDestPath, "ReferenceGetSortKeysByPartitionKey.go")
	tp.CreateFileFromTemplate(templ, nil, referenceSortKeyPath)

	templ = tp.ParseOrExitOnError("templates/type_spec/add_many_to_may.template.go.tmpl")
	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "reference", "AddToManyToMany"))
	os.MkdirAll(generationDestPath, 0777)
	addManyToManyPath := filepath.Join(generationDestPath, "ReferenceAddToManyToMany.go")
	tp.CreateFileFromTemplate(templ, nil, addManyToManyPath)
}

func GenerateCustomConstructorsHandlers(path string, customCtor []parser.CustomCtorDefinition) {
	var handlerDir string
	var customCtorFileName string
	templ := tp.ParseOrExitOnError("templates/type_spec/custom_constructor_template.go.tmpl")
	generationDestPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "custom-constructors"))

	for _, c := range customCtor {
		customCtorFileName = "New" + c.TypeName
		handlerDir = filepath.Join(generationDestPath, customCtorFileName)
		os.MkdirAll(handlerDir, 0777)
		path = filepath.Join(handlerDir, customCtorFileName+".go")
		tp.CreateFileFromTemplate(templ, c, path)
		tp.RunGoimportsOnFile(path)
	}
}
