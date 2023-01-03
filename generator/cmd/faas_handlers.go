package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Astenna/Nubes/generator/database"
	"github.com/Astenna/Nubes/generator/parser"
	tp "github.com/Astenna/Nubes/generator/template_parser"
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
		dbInit, _ := cmd.Flags().GetBool("dbInit")
		generateDeploymentFiles, _ := cmd.Flags().GetBool("deplFiles")

		typesPath = tp.MakePathAbosoluteOrExitOnError(typesPath)
		repositoriesPath = tp.MakePathAbosoluteOrExitOnError(repositoriesPath)

		isNobjectType, nobjectsImportPath := parser.GetNobjectsDefinedInPack(typesPath, moduleName)
		stateChangingFuncs := parser.ParseStateChangingHandlers(typesPath, nobjectsImportPath, isNobjectType)
		customRepoFuncs, defaultRepoFuncs := parser.ParseRepoHandlers(repositoriesPath, nobjectsImportPath, isNobjectType)

		GenerateStateChangingHandlers(generationDestination, stateChangingFuncs)
		GenerateRepositoriesHandlers(generationDestination, customRepoFuncs, defaultRepoFuncs)

		if generateDeploymentFiles {
			serviceName := lastString(strings.Split(moduleName, "/"))
			serverlessInput := ServerlessTemplateInput{ServiceName: serviceName, DefaultRepos: defaultRepoFuncs, CustomRepos: customRepoFuncs, StateFuncs: stateChangingFuncs}
			GenerateDeploymentFiles(generationDestination, serverlessInput)
		}

		if dbInit {
			database.CreateTypeTables(isNobjectType)
		}
	},
}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}

func init() {
	rootCmd.AddCommand(handlersCmd)

	var typesPath string
	var repositoriesPath string
	var handlersPath string
	var moduleName string
	var dbInit bool
	var generateDeploymentFiles bool

	handlersCmd.Flags().StringVarP(&typesPath, "types", "t", ".", "path to directory with types")
	handlersCmd.Flags().StringVarP(&repositoriesPath, "repositories", "r", ".", "path to directory with repositories")
	handlersCmd.Flags().StringVarP(&handlersPath, "output", "o", ".", "path where directory with handlers will be created")
	handlersCmd.Flags().StringVarP(&moduleName, "module", "m", "MISSING_MODULE_NAME", "module name of the source project")
	handlersCmd.Flags().BoolVarP(&dbInit, "dbInit", "i", false, "boolean, indicates whether database should be initialized by creation of tables based on type names")
	handlersCmd.Flags().BoolVarP(&generateDeploymentFiles, "deplFiles", "g", true, "boolean, indicates whether deployment files for AWS lambdas are to be created")

	cmd.Execute()
}

type ServerlessTemplateInput struct {
	ServiceName  string
	DefaultRepos []parser.DefaultRepoHandler
	CustomRepos  []parser.CustomRepoHandler
	StateFuncs   []parser.StateChangingHandler
}

func GenerateDeploymentFiles(path string, templateInput ServerlessTemplateInput) {
	serverlessTempl := tp.ParseOrExitOnError("templates/handlers/deployment/serverless.yml.tmpl")
	fileName := filepath.Join(tp.MakePathAbosoluteOrExitOnError(path), "serverless.yml")
	tp.CreateFileFromTemplate(serverlessTempl, templateInput, fileName)

	buildScriptTempl := tp.ParseOrExitOnError("templates/handlers/deployment/build_handlers.sh.tmpl")
	fileName = filepath.Join(tp.MakePathAbosoluteOrExitOnError(path), "build_handlers.sh")
	tp.CreateFileFromTemplate(buildScriptTempl, nil, fileName)

	dockerfileTempl := tp.ParseOrExitOnError("templates/handlers/deployment/Dockerfile.tmpl")
	fileName = filepath.Join(tp.MakePathAbosoluteOrExitOnError(path), "Dockerfile")
	tp.CreateFileFromTemplate(dockerfileTempl, nil, fileName)
}

func GenerateRepositoriesHandlers(path string, customFuncs []parser.CustomRepoHandler, defaultFuncs []parser.DefaultRepoHandler) {
	var fileName string
	var handlerDir string
	var operationTypeCombined string
	templ := tp.ParseOrExitOnError("templates/handlers/custom_repo_template.go.tmpl")
	repositoriesDirectoryPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "repositories"))

	for _, f := range customFuncs {
		operationTypeCombined = f.OperationName + f.TypeName
		handlerDir = filepath.Join(repositoriesDirectoryPath, operationTypeCombined)
		os.MkdirAll(handlerDir, 0777)
		fileName = filepath.Join(handlerDir, operationTypeCombined+".go")
		tp.CreateFileFromTemplate(templ, f, fileName)
	}

	getTempl := tp.ParseOrExitOnError("templates/handlers/get_repo_template.go.tmpl")
	createTempl := tp.ParseOrExitOnError("templates/handlers/create_repo_template.go.tmpl")
	deleteTempl := tp.ParseOrExitOnError("templates/handlers/delete_repo_template.go.tmpl")
	updateTempl := tp.ParseOrExitOnError("templates/handlers/update_repo_template.go.tmpl")

	var tmpl template.Template
	for _, f := range defaultFuncs {
		switch {
		case f.OperationName == parser.GetPrefix:
			tmpl = getTempl
		case f.OperationName == parser.CreatePrefix:
			tmpl = createTempl
		case f.OperationName == parser.DeletePrefix:
			tmpl = deleteTempl
		case f.OperationName == parser.UpdatePrefix:
			tmpl = updateTempl
		}

		operationTypeCombined = f.OperationName + f.TypeName
		handlerDir = filepath.Join(repositoriesDirectoryPath, operationTypeCombined)
		os.MkdirAll(handlerDir, 0777)
		fileName = filepath.Join(handlerDir, operationTypeCombined+".go")
		tp.CreateFileFromTemplate(tmpl, f, fileName)
		tp.RunGoimportsOnFile(fileName)
	}
}

func GenerateStateChangingHandlers(path string, functions []parser.StateChangingHandler) {
	var handlerDir string
	var ownerHandlerNameCombined string
	templ := tp.ParseOrExitOnError("templates/handlers/state_changing_template.go.tmpl")
	generationDestPath := tp.MakePathAbosoluteOrExitOnError(filepath.Join(path, "generated", "state-changes"))

	for _, f := range functions {
		ownerHandlerNameCombined = f.OwnerType + f.HandlerNameWithoutSuffix
		handlerDir = filepath.Join(generationDestPath, ownerHandlerNameCombined)
		os.MkdirAll(handlerDir, 0777)
		filepath := filepath.Join(handlerDir, ownerHandlerNameCombined+".go")
		tp.CreateFileFromTemplate(templ, f, filepath)
		tp.RunGoimportsOnFile(filepath)
	}
}
