package cmd

import (
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

		parsedPackage := parser.GetPackageTypes(typesPath, moduleName)
		stateChangingFuncs := parser.ParseStateChangingHandlers(typesPath, parsedPackage)
		parser.AddDBOperationsToMethods(typesPath, parsedPackage)

		GenerateStateChangingHandlers(generationDestination, stateChangingFuncs)
		GenerateGetAndSetFieldHandlers(generationDestination)

		if generateDeploymentFiles {
			serviceName := lastString(strings.Split(moduleName, "/"))
			serverlessInput := ServerlessTemplateInput{ServiceName: serviceName, StateFuncs: stateChangingFuncs}
			GenerateDeploymentFiles(generationDestination, serverlessInput)
		}

		if dbInit {
			database.CreateTypeTables(parsedPackage)
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
}

func GenerateDeploymentFiles(path string, templateInput ServerlessTemplateInput) {
	serverlessTempl := tp.ParseOrExitOnError("templates/ssf_spec/deployment/serverless.yml.tmpl")
	fileName := filepath.Join(tp.MakePathAbosoluteOrExitOnError(path), "serverless.yml")
	tp.CreateFileFromTemplate(serverlessTempl, templateInput, fileName)

	buildScriptTempl := tp.ParseOrExitOnError("templates/ssf_spec/deployment/build_handlers.sh.tmpl")
	fileName = filepath.Join(tp.MakePathAbosoluteOrExitOnError(path), "build_handlers.sh")
	tp.CreateFileFromTemplate(buildScriptTempl, nil, fileName)
}

func GenerateStateChangingHandlers(path string, functions []parser.StateChangingHandler) {
	var handlerDir string
	var ownerHandlerNameCombined string
	templ := tp.ParseOrExitOnError("templates/ssf_spec/state_changing_template.go.tmpl")
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

func GenerateGetAndSetFieldHandlers(path string) {
	templ := tp.ParseOrExitOnError("templates/ssf_spec/get_field_template.go.tmpl")
	generationDestPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "GetField"))
	os.MkdirAll(generationDestPath, 0777)
	getPath := filepath.Join(generationDestPath, "GetField.go")
	tp.CreateFileFromTemplate(templ, nil, getPath)

	templ = tp.ParseOrExitOnError("templates/ssf_spec/set_field_template.go.tmpl")
	generationDestPath = tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "generics", "SetField"))
	os.MkdirAll(generationDestPath, 0777)
	setPath := filepath.Join(generationDestPath, "SetField.go")
	tp.CreateFileFromTemplate(templ, nil, setPath)
}
