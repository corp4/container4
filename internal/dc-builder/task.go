package builder

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/corp4/container4/pkg/devcontainer"
	"github.com/corp4/container4/pkg/git"
	"github.com/corp4/container4/pkg/websocket"
)

// Return a task from a base64 encoded json string
func getTaskFromString(taskInfoStr string) (Task, error) {
	if taskInfoStr == "" {
		return Task{}, fmt.Errorf("task info data is empty")
	}

	// Base 64 decode the string
	jsonStr, err := base64.StdEncoding.DecodeString(taskInfoStr)
	if err != nil {
		return Task{}, err
	}

	log.Debugln("Task info string:", string(jsonStr))

	// Unmarshal the json string
	var task Task
	err = json.Unmarshal(jsonStr, &task)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func createWorkspacesFolder() (string, error) {
	// Create the workspace directory
	workspacesPath := "/tmp/workspaces"
	err := os.MkdirAll(workspacesPath, 0755)
	if err != nil {
		return "", err
	}

	return workspacesPath, nil
}

func RunTask(taskInfoStr string, ws *websocket.Websocket) {
	// Get the task from the taskInfo string
	task, err := getTaskFromString(taskInfoStr)
	if err != nil {
		log.Errorln("Error getting task from taskInfo:", err)
		return
	}

	log.Infoln("Running task: ", task.TaskName)

	// Create the workspace directory
	workspacesPath, err := createWorkspacesFolder()
	if err != nil {
		log.Errorln("Error creating workspace directory:", err)
		return
	}

	workspacePath := fmt.Sprintf("%s/%s", workspacesPath, task.Workspace.Name)

	// Remove folder if it already exists
	err = os.RemoveAll(workspacePath)
	if err != nil {
		log.Errorln("Error removing workspace directory:", err)
		return
	}

	// Clone the repo in the workspace directory
	err = git.CloneRepo(task.Repo, task.Provider, workspacePath)
	if err != nil {
		log.Errorln("Error cloning repo:", err)
		return
	}

	log.Infoln("Successfully cloned repo")

	dc := devcontainer.NewDevcontainer(workspacePath, "", ws)

	// Get and merge the devcontainer config with the service one
	configPath, err := mergeConfigDevcontainer(task, dc)
	if err != nil {
		log.Errorln("Error merging devcontainer config:", err)
		return
	}

	dc.OverridedConfig = configPath
	log.Infoln("Successfully merged devcontainer config")

	// Create the devcontainer
	log.Infoln("Creating devcontainer with config:", configPath)
	err = setupDevcontainer(dc)
	if err != nil {
		log.Errorln("Error creating devcontainer:", err)
		return
	}

	log.Infoln("Successfully created devcontainer")
}
