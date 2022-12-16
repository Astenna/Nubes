package cmd

import (
	"github.com/Astenna/Thesis_PoC/generator/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/cmd"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generates client project",
	Long:  `Generates client project based on types and repositories.`,

	Run: func(cmd *cobra.Command, args []string) {
		typesPath, _ := cmd.Flags().GetString("types")
		repositoriesPath, _ := cmd.Flags().GetString("repositories")
		//output, _ := cmd.Flags().GetString("output")
		_ = repositoriesPath

		parser.PrepareTypesFiles(MakePathAbosoluteOrExitOnError(typesPath))
		//types := parser.PrepareTypes(MakePathAbosoluteOrExitOnError(typesPath))
		//templ, _ := template.ParseFiles("type_template.go.tmpl")

		// outputDirectoryPath := MakePathAbosoluteOrExitOnError(filepath.Join(output, "output_testing"))
		// os.MkdirAll(outputDirectoryPath, 0777)

		// for i, f := range types {
		// 	file, err := os.Create(filepath.Join(outputDirectoryPath, "test"+strconv.Itoa(i)+".go"))
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	err = templ.Execute(file, f)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	defer file.Close()
		// }
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	var typesPath string
	var repositoriesPath string
	var outputPath string

	clientCmd.Flags().StringVarP(&typesPath, "types", "t", ".", "path to directory with types")
	clientCmd.Flags().StringVarP(&repositoriesPath, "repositories", "r", ".", "path to directory with repositories")
	clientCmd.Flags().StringVarP(&outputPath, "output", "o", ".", "path where directory with client project will be created")

	cmd.Execute()
}
