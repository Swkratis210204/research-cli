package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base "research" command
// every other command (arxiv, search, etc.) will be added to this
var rootCmd = &cobra.Command{
	Use:   "research",                                                                         // the command name
	Short: "A personal research assistant CLI",                                                // shown in --help
	Long:  "research is a CLI tool for fetching, summarizing and organizing academic papers.", // shown in research --help
}

// Execute is called by main.go to start the CLI
func Execute() {
	// run cobra, which reads what the user typed and routes to the right command
	// := runs rootCmd.Execute() and stores the result in err
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
