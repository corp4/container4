package builder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/corp4/container4/pkg/devcontainer"
	log "github.com/sirupsen/logrus"
)

func buildDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":  "Flashenv Default Environment",
		"image": "mcr.microsoft.com/devcontainers/base:jammy",
	}
}

func buildServiceConfig(task Task) map[string]interface{} {
	return map[string]interface{}{
		"runArgs": []string{
			"--label", "FlashEnvVersion=0.0.1",
			"--hostname", task.Workspace.Name,
			"--add-host", fmt.Sprintf("%s:127.0.0.1", task.Workspace.Name),
			"--cap-add=sys_nice",
			"--network", "host",
		},
		"containerEnv": map[string]string{
			"FLASHENV":                   "true",
			"FLASHENV_VERSION":           "0.0.1",
			"FLASHENV_WORKSPACE_NAME":    task.Workspace.Name,
			"FLASHENV_WORKSPACE_USERID":  task.Workspace.UserId,
			"FLASHENV_REPO_URL":          task.Repo.Url,
			"FLASHENV_REPO_BRANCH":       task.Repo.Branch,
			"FLASHENV_REPO_COMMIT":       task.Repo.Commit,
			"FLASHENV_PROVIDER_NAME":     task.Provider.Name,
			"FLASHENV_PROVIDER_TOKEN":    task.Provider.Token,
			"FLASHENV_PROVIDER_USERNAME": task.Provider.Username,
		},
		"remoteEnv": map[string]string{
			"FLASHENV":                   "true",
			"FLASHENV_VERSION":           "0.0.1",
			"FLASHENV_WORKSPACE_NAME":    task.Workspace.Name,
			"FLASHENV_WORKSPACE_USERID":  task.Workspace.UserId,
			"FLASHENV_REPO_URL":          task.Repo.Url,
			"FLASHENV_REPO_BRANCH":       task.Repo.Branch,
			"FLASHENV_REPO_COMMIT":       task.Repo.Commit,
			"FLASHENV_PROVIDER_NAME":     task.Provider.Name,
			"FLASHENV_PROVIDER_TOKEN":    task.Provider.Token,
			"FLASHENV_PROVIDER_USERNAME": task.Provider.Username,
		},
		//"remoteUser": "vscode",
	}
}

// Merge the service config with the devcontainer config
// Overwrite the values in the devcontainer config with the service ones
func mergeConfig(config map[string]interface{}, serviceConfig map[string]interface{}) (map[string]interface{}, error) {
	for key, value := range serviceConfig {
		// Check if the key is present in the config
		if _, ok := config[key]; !ok {
			// Key not present, add it
			config[key] = value
		} else {
			// Key present, check if it's a map
			switch config[key].(type) {
			case map[string]interface{}:
				// It's a map, merge the two maps
				// Check if the value is a map
				switch value.(type) {
				case map[string]interface{}:
					// It's a map, merge the two maps
					newMap, err := mergeConfig(config[key].(map[string]interface{}), value.(map[string]interface{}))
					if err != nil {
						return nil, err
					}
					config[key] = newMap
				default:
					// Error, can't merge a map with a non-map
					return nil, fmt.Errorf("can't merge a map with a non-map")
				}
			default:
				// It's not a map, overwrite the value with the service one
				config[key] = value
			}
		}
	}

	return config, nil
}

func saveConfigDevcontainer(config map[string]interface{}) (string, error) {
	// Convert the config to json
	configJson, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	log.Debugln("Merged config:")
	log.Debugln(configJson)

	// Get the folder where to save the config
	folderConfig := filepath.Join(os.Getenv("HOME"), "builder")

	// Create the folder if it doesn't exist
	log.Debugln("Creating folder:", folderConfig)
	err = os.MkdirAll(folderConfig, 0755)
	if err != nil {
		return "", err
	}

	// Save the config as json in the home folder
	configPath := filepath.Join(folderConfig, "merged_config.json")

	// Remove file if it exists
	log.Debugln("Removing file:", configPath)
	err = os.Remove(configPath)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	log.Debugln("Creating file:", configPath)
	err = os.WriteFile(configPath, configJson, 0644)
	if err != nil {
		return "", err
	}

	return configPath, nil
}

func mergeConfigDevcontainer(task Task, dc *devcontainer.Devcontainer) (string, error) {
	// Check if the workspace folder is a devcontainer with a valid config file
	config, err := dc.GetConfig()
	if err != nil {
		log.Warnln("Error reading devcontainer config:", err)

		// No valid config, set the default one
		config = buildDefaultConfig()
	}

	// Merge the config with the service one
	serviceConfig := buildServiceConfig(task)
	config, err = mergeConfig(config, serviceConfig)
	if err != nil {
		return "", err
	}

	// Save the config as json in the home folder
	return saveConfigDevcontainer(config)
}
