package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/cmd"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generates client project",
	Long:  `Generates client project based on types and repositories.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called!")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	cmd.Execute()
}
