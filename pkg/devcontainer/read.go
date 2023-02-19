package devcontainer

import (
	"encoding/json"

	"github.com/corp4/container4/pkg/execute"
	log "github.com/sirupsen/logrus"
)

func (dc *Devcontainer) GetConfig() (map[string]interface{}, error) {
	// Prepare arguments for the devcontainer command
	args := []string{"--workspace-folder", dc.WorkspaceFolder}

	devcontainerCmd := append([]string{"read-configuration"}, args...)
	cmd, stdout, stderr, err := execute.AsyncExecute(append([]string{"devcontainer"}, devcontainerCmd...))
	if err != nil {
		return nil, err
	}

	// Read the output
	out, err := execute.SyncReadCommandResult(cmd, stdout, stderr)
	if err != nil {
		return nil, err
	}

	// Parse the output
	var config map[string]interface{}
	err = json.Unmarshal([]byte(out), &config)
	if err != nil {
		return nil, err
	}

	// Get the "configuration" from the map
	configuration, ok := config["configuration"].(map[string]interface{})
	if !ok {
		return nil, err
	}

	// Remove the "configFilePath" from the configuration
	delete(configuration, "configFilePath")

	log.Debugln("Devcontainer configuration:", configuration)

	return configuration, nil
}
