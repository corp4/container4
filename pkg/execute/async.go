package execute

import (
	"io"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func AsyncExecute(args []string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	// Create the command
	cmd := exec.Command(args[0], args[1:]...)

	// Get the stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorln("Error getting stdout pipe:", err)
		return nil, nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Errorln("Error getting stderr pipe:", err)
		return nil, nil, nil, err
	}

	// Start the command
	log.Infoln("Executing:", args)
	err = cmd.Start()
	log.Debugln("Command started: ", cmd.Process.Pid)
	if err != nil {
		log.Errorln("Error executing command:", err)
		return nil, nil, nil, err
	}

	return cmd, stdout, stderr, nil
}
