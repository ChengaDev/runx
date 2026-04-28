package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "runx",
	Short: "Run any saved project from anywhere on your system",
	Long: `runx lets you save projects with a name, path, and run command,
then launch them from anywhere without navigating to their directory.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(editCmd)
}
