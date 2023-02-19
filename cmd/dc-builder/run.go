package main

import (
	builder "github.com/corp4/container4/internal/dc-builder"
	"github.com/corp4/container4/pkg/websocket"
	"github.com/spf13/cobra"
)

var (
	// Information about the task to run as json base64 encoded string
	taskInfo string

	// Websockets information
	host string
	port int
)

// Run the building task
var runCmd = &cobra.Command{
	Use:        "run",
	Short:      "Run the building task",
	Long:       "Setup the instance with the given task and run it",
	SuggestFor: []string{"serve", "start"},
	Run: func(cmd *cobra.Command, args []string) {
		// Create web socket
		ws := websocket.CreateWebsocket(host, port)
		ws.Start()
		builder.RunTask(taskInfo, ws)
	},
}

// Initialize the run command
func init() {
	rootCmd.AddCommand(runCmd)

	// Add must have flag named userInfos that is a base64 string
	runCmd.Flags().StringVarP(&taskInfo, "taskInfo", "t", "", "Task information as json base64 encoded string")

	// Set websocket flags
	runCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Websocket host")
	runCmd.Flags().IntVarP(&port, "port", "p", 8080, "Websocket port")

	runCmd.MarkFlagRequired("taskInfo")
}
