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
	Use:     "dc-supervisor",
	Short:   "RPC server to control devcontainers",
	Long: `dc-supervisor is a RPC over HTTP server designed to provide efficient and
convenient access to data and actions inside devcontainers.`,
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
