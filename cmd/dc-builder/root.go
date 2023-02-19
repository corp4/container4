package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Sets verbose mode
	verbose bool
)

// Base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.0.1",
	Use:     "dc-builder",
	Short:   "Setup and build devcontainers from repository of a user",
	Long:    "Setup and build devcontainers from repository of a user",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Info("Debug mode enabled")
		}
	},
}

// Initialize the root command
func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

// Starts cobra
func main() {
	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
