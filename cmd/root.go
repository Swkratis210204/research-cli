package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "research",
	Short: "A personal research assistant CLI",
	Long: `research is a CLI tool for fetching, summarizing and organizing academic papers.

Docs & issues: https://github.com/Swkratis210204/research-cli`,
}

// Execute is called by main.go to start the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init hides the completion command from the help menu
func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
