package devcontainer

import (
	"github.com/corp4/container4/pkg/execute"
	"github.com/corp4/container4/pkg/websocket"
	log "github.com/sirupsen/logrus"
)

type Devcontainer struct {
	WorkspaceFolder string   // Path to the workspace folder
	OverridedConfig string   // Path to the overrided config file
	Args            []string // Default arguments for the devcontainer command

	// The websocket to redirect the output to
	Websocket *websocket.Websocket
}

// Create a new devcontainer
func NewDevcontainer(workspaceFolder string, overridedConfigFile string, ws *websocket.Websocket) *Devcontainer {
	return &Devcontainer{
		WorkspaceFolder: workspaceFolder,
		OverridedConfig: overridedConfigFile,
		Websocket:       ws,
		Args: []string{
			"--id-label", "Type=flashenv",
			"--log-level", "trace",
			"--update-remote-user-uid-default", "never",
			"--mount-workspace-git-root", "false",
			"--default-user-env-probe", "loginInteractiveShell",
		},
	}
}

func (dc *Devcontainer) Infoln(msg string) {
	log.Infoln(msg)
	if dc.Websocket != nil {
		dc.Websocket.Send([]byte(msg))
	}
}

// Run the specified command of the devcontainer-cli tool with the default arguments and the given arguments
func (dc *Devcontainer) RunCommand(cmdName string, logName string, args []string) error {
	// Prepare arguments for the devcontainer command
	devcontainerCmd := append([]string{cmdName}, dc.Args...)
	devcontainerCmd = append(devcontainerCmd, args...)

	dc.Infoln("=== Devcontainer " + cmdName + ": " + logName + " Command ===")
	cmd, stdout, stderr, err := execute.AsyncExecute(append([]string{"devcontainer"}, devcontainerCmd...))
	if err != nil {
		return err
	}

	// Redirect the output to the websocket if needed
	if dc.Websocket != nil {
		err = dc.Websocket.RedirectCommandOutput(cmd, stdout, stderr)
		if err != nil {
			return err
		}
	}

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
