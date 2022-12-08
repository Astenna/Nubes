package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/template"

	"github.com/Astenna/Thesis_PoC/generator/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/cmd"
)

var handlersCmd = &cobra.Command{
	Use:   "handlers",
	Short: "Generates handlers for AWS lambda deployment",
	Long:  `Generates handlers for AWS lambda deployment based on types and repositories indicated by the path`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("handlers called!")

		typesPath, _ := cmd.Flags().GetString("types")
		repositoriesPath, _ := cmd.Flags().GetString("repositories")
		fmt.Println("typesPath: ", typesPath)
		fmt.Println("repositoriesPath: ", repositoriesPath)

		subPackage := "C:\\Users\\marek\\OneDrive\\master-thesis\\Thesis_PoC\\faas\\types"
		functions := parser.PrepareHandlerFunctions(subPackage)
		//parser.ParseTypes(subPackage)

		templ, _ := template.ParseFiles("handler_template.tmpl")

		for i, f := range functions {
			file, _ := os.Create("test" + strconv.Itoa(i) + ".go")

			err := templ.Execute(file, f)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()
		}

	},
}

func init() {
	rootCmd.AddCommand(handlersCmd)

	var typesPath string
	var repositoriesPath string
	var handlersPath string

	handlersCmd.Flags().StringVarP(&typesPath, "types", "t", ".", "path to directory with types")
	handlersCmd.Flags().StringVarP(&repositoriesPath, "repositories", "r", ".", "path to directory with repositories")
	handlersCmd.Flags().StringVarP(&handlersPath, "output", "o", ".", "path where directory with handlers will be created")

	cmd.Execute()
}
