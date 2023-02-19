package devcontainer

import (
	"encoding/json"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func (dc *Devcontainer) GetConfig() (map[string]interface{}, error) {
	// Prepare arguments for the devcontainer command
	args := []string{"read-configuration", "--workspace-folder", dc.WorkspaceFolder}

	// Read the output
	out, err := exec.Command("devcontainer", args...).Output()
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
