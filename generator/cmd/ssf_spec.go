package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Astenna/Nubes/generator/database"
	"github.com/Astenna/Nubes/generator/parser"
	tp "github.com/Astenna/Nubes/generator/template"
	typespec "github.com/Astenna/Nubes/generator/template/type_spec"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/cmd"
)

var ssfSpecCmd = &cobra.Command{
	Use:   "handlers",
	Short: "Generates handlers' definitions for AWS lambda deployment",
	Long:  `Generates handlers' definitions for AWS lambda deployment based on types indicated by the path`,

	Run: func(cmd *cobra.Command, _ []string) {
		typesPath, _ := cmd.Flags().GetString("types")
		generationDestination, _ := cmd.Flags().GetString("output")
		moduleName, _ := cmd.Flags().GetString("module")
		dbInit, _ := cmd.Flags().GetBool("dbInit")
		generateDeploymentFilesOn, _ := cmd.Flags().GetBool("deplFiles")

		typesPath = tp.MakePathAbosoluteOrExitOnError(typesPath)

		typeSpecParser, err := parser.NewTypeSpecParser(typesPath)
		if err != nil {
			fmt.Println("Fatal error occurred initialising type spec parser: %w", err)
			os.Exit(1)
		}
		typeSpecParser.Run(moduleName)

		generateStateChangingHandlers(generationDestination, typeSpecParser.Handlers)
		generateGenericHandlers(generationDestination, typeSpecParser.Output)
		generateCustomConstructorsHandlers(generationDestination, typeSpecParser.CustomCtors)

		if generateDeploymentFilesOn {
			serviceName := lastElem(strings.Split(moduleName, "/"))
			serverlessInput := ServerlessTemplateInput{
				ServiceName:   serviceName,
				StateFuncs:    typeSpecParser.Handlers,
				CustomCtors:   typeSpecParser.CustomCtors,
				ManyToManyRel: len(typeSpecParser.Output.ManyToManyRelationships) > 0,
			}
			generateDeploymentFiles(generationDestination, serverlessInput)
		}

		if dbInit {
			database.CreateTypeTables(typeSpecParser.Output)
		}
	},
}

func init() {
	rootCmd.AddCommand(ssfSpecCmd)

	var typesPath string
	var handlersPath string
	var moduleName string
	var dbInit bool
	var generateDeploymentFiles bool

	ssfSpecCmd.Flags().StringVarP(&typesPath, "types", "t", ".", "path to package with types definitions")
	ssfSpecCmd.Flags().StringVarP(&handlersPath, "output", "o", ".", "path where directory with handlers will be created")
	ssfSpecCmd.Flags().StringVarP(&moduleName, "module", "m", "MISSING_MODULE_NAME", "module name of the source project")
	ssfSpecCmd.Flags().BoolVarP(&dbInit, "dbInit", "i", false, "boolean, indicates whether database tables should be initialized")
	ssfSpecCmd.Flags().BoolVarP(&generateDeploymentFiles, "deplFiles", "g", true, "boolean, indicates whether deployment files for AWS lambdas are to be created")

	cmd.Execute()
}

type ServerlessTemplateInput struct {
	ServiceName   string
	StateFuncs    []parser.StateChangingHandler
	CustomCtors   []parser.CustomCtorDefinition
	ManyToManyRel bool
}

func generateDeploymentFiles(path string, templateInput ServerlessTemplateInput) {
	fileName := filepath.Join(tp.MakePathAbosoluteOrExitOnError(path), "serverless.yml")
	tp.CreateFile("template/type_spec/deployment/serverless.yml.tmpl", templateInput, fileName)

	fileName = filepath.Join(tp.MakePathAbosoluteOrExitOnError(path), "build_handlers.sh")
	tp.CreateFile("template/type_spec/deployment/build_handlers.sh.tmpl", nil, fileName)
}

func generateStateChangingHandlers(path string, functions []parser.StateChangingHandler) {
	var handlerDir string
	var ownerHandlerNameCombined string
	generationDestPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "state-changes"))

	for _, f := range functions {
		ownerHandlerNameCombined = f.ReceiverType + f.MethodName
		handlerDir = filepath.Join(generationDestPath, ownerHandlerNameCombined)
		os.MkdirAll(handlerDir, 0777)
		path = filepath.Join(handlerDir, ownerHandlerNameCombined+".go")
		tp.CreateFile("template/type_spec/state_changing_template.go.tmpl", f, path)
		tp.RunGoimportsOnFile(path)
	}
}

func generateGenericHandlers(path string, parsedPkg parser.ParsedPackage) {
	generationDestPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "GetBatch"))
	os.MkdirAll(generationDestPath, 0777)
	getBatch := filepath.Join(generationDestPath, "GetBatch.go")
	tp.CreateFile("template/type_spec/get_batch.go.tmpl", parsedPkg.TypesWithCustomId, getBatch)

	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "GetState"))
	os.MkdirAll(generationDestPath, 0777)
	getPath := filepath.Join(generationDestPath, "GetState.go")
	tp.CreateFile("template/type_spec/get_state.go.tmpl", parsedPkg.TypesWithCustomId, getPath)

	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "SetField"))
	os.MkdirAll(generationDestPath, 0777)
	setPath := filepath.Join(generationDestPath, "SetField.go")
	tp.CreateFile("template/type_spec/set_field_template.go.tmpl", nil, setPath)

	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "Load"))
	os.MkdirAll(generationDestPath, 0777)
	loadPath := filepath.Join(generationDestPath, "Load.go")
	tp.CreateFile("template/type_spec/load_template.go.tmpl", nil, loadPath)

	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "Export"))
	os.MkdirAll(generationDestPath, 0777)
	exportPath := filepath.Join(generationDestPath, "Export.go")
	input := typespec.ExportTemplateInput{IsNobjectInOrginalPackage: parsedPkg.IsNobjectInOrginalPackage,
		OrginalPackageAlias: parser.OrginalPackageAlias, OrginalPackage: parsedPkg.ImportPath,
		TypesWithCustomExport: parsedPkg.TypesWithCustomExport,
	}
	tp.CreateFile("template/type_spec/export_template.go.tmpl", input, exportPath)
	tp.RunGoimportsOnFile(generationDestPath)

	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "Delete"))
	os.MkdirAll(generationDestPath, 0777)
	deletePath := filepath.Join(generationDestPath, "Delete.go")
	deleteTemplInput := typespec.DeleteTemplateInput{OrginalPackageAlias: parser.OrginalPackageAlias,
		OrginalPackage:        parsedPkg.ImportPath,
		TypesWithCustomDelete: parsedPkg.TypesWithCustomDelete,
	}
	tp.CreateFile("template/type_spec/delete_template.go.tmpl", deleteTemplInput, deletePath)
	tp.RunGoimportsOnFile(generationDestPath)

	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "reference", "Get"))
	os.MkdirAll(generationDestPath, 0777)
	referenceIndexPath := filepath.Join(generationDestPath, "ReferenceGet.go")
	tp.CreateFile("template/type_spec/reference_get.go.tmpl", nil, referenceIndexPath)

	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "reference", "GetIds"))
	os.MkdirAll(generationDestPath, 0777)
	referenceGetIdsPath := filepath.Join(generationDestPath, "ReferenceGetIds.go")
	tp.CreateFile("template/type_spec/reference_get_ids.go.tmpl", nil, referenceGetIdsPath)

	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "reference", "GetStubs"))
	os.MkdirAll(generationDestPath, 0777)
	referenceGetStubsPath := filepath.Join(generationDestPath, "ReferenceGetStubs.go")
	tp.CreateFile("template/type_spec/reference_get_stubs.go.tmpl", nil, referenceGetStubsPath)

	if len(parsedPkg.ManyToManyRelationships) > 0 {

		generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "reference", "AddToManyToMany"))
		os.MkdirAll(generationDestPath, 0777)
		addManyToManyPath := filepath.Join(generationDestPath, "ReferenceAddToManyToMany.go")
		tp.CreateFile("template/type_spec/add_many_to_many.go.tmpl", nil, addManyToManyPath)

		generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "reference", "DeleteFromManyToMany"))
		os.MkdirAll(generationDestPath, 0777)
		deleteFromManyPath := filepath.Join(generationDestPath, "ReferenceDeleteFromManyToMany.go")
		tp.CreateFile("template/type_spec/delete_from_many_to_many.go.tmpl", nil, deleteFromManyPath)
	}
}

func generateCustomConstructorsHandlers(path string, customCtor []parser.CustomCtorDefinition) {
	var handlerDir string
	var customCtorFileName string
	generationDestPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "custom-constructors"))

	for _, c := range customCtor {
		customCtorFileName = "New" + c.TypeName
		handlerDir = filepath.Join(generationDestPath, customCtorFileName)
		os.MkdirAll(handlerDir, 0777)
		path = filepath.Join(handlerDir, customCtorFileName+".go")
		tp.CreateFile("template/type_spec/custom_constructor.go.tmpl", c, path)
		tp.RunGoimportsOnFile(path)
	}
}

func lastElem(ss []string) string {
	return ss[len(ss)-1]
}
