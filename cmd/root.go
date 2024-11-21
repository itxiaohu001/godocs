package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godoc-gen",
	Short: "A documentation generator for Go projects",
	Long: `A documentation generator that creates comprehensive Markdown documentation
for Go projects, including struct definitions, field information, and comments.`,
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
