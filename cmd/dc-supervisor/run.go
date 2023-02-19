package main

import (
	"fmt"
	"net/http"
	"net/rpc"

	supervisor "github.com/corp4/container4/internal/dc-supervisor"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	// Host to listen on
	host string

	// Port to listen on
	port string
)

// Run the RPC server
var runCmd = &cobra.Command{
	Use:        "run",
	Short:      "Run the RPC server",
	Long:       `Initialize the RPC server and listen for incoming connections`,
	SuggestFor: []string{"serve", "start"},
	Run: func(cmd *cobra.Command, args []string) {

		// Register RPC services
		var ssh supervisor.SSH
		if err := rpc.Register(&ssh); err != nil {
			log.Fatal(err)
		}
		log.Debug("Registered SSH service")

		var supervisor supervisor.Supervisor
		if err := rpc.Register(&supervisor); err != nil {
			log.Fatal(err)
		}
		log.Debug("Registered Supervisor service")

		rpc.HandleHTTP()

		// Start listening for incoming connections
		log.Infof("Listening on %s:%s", host, port)
		if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil); err != nil {
			log.Fatal(err)
		}
	},
}

// Initialize the run command
func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "Host to listen on")
	runCmd.Flags().StringVarP(&port, "port", "p", "38156", "Port to listen on")
}
