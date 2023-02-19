package builder

import "github.com/corp4/container4/pkg/devcontainer"

// Run the 3 steps for the full creation of a devcontainer
func setupDevcontainer(dc *devcontainer.Devcontainer) error {
	// Create arguments for the devcontainer command
	baseArgs := []string{"--workspace-folder", dc.WorkspaceFolder, "--override-config", dc.OverridedConfig}

	// Step 1: Create the devcontainer
	args := append(baseArgs, "--skip-non-blocking-commands", "--skip-post-create")
	err := dc.RunCommand("up", "Create", args)
	if err != nil {
		return err
	}

	// Step 2: Run the blocking commands
	args = append(baseArgs, "--skip-non-blocking-commands", "--expect-existing-container")
	err = dc.RunCommand("up", "Run Blocking", args)
	if err != nil {
		return err
	}

	// Step 3: Run the devcontainer
	args = append(baseArgs, "--skip-post-attach", "--expect-existing-container")
	err = dc.RunCommand("up", "Run Command", args)
	if err != nil {
		return err
	}

	return nil
}
