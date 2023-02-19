package execute

import (
	"io"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Run a command and wait for it to finish while storing the output in a string
func SyncReadCommandResult(cmd *exec.Cmd, stdout io.ReadCloser, stderr io.ReadCloser) (string, error) {
	var out string
	buf := make([]byte, 1024)

	// Wait for the command to finish
	err := cmd.Wait()

	// Read stdout even if there is an error
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			log.Debugln("End of output stdout: ", err)
			break
		}

		log.Debugln("Read: ", string(buf[:n]))
		out += string(buf[:n])
	}

	// If there is an error, read stderr and return the error
	if err != nil {
		log.Debugln("Command finished with error: ", err)

		// Read stderr
		buf = make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				log.Debugln("End of output stderr: ", err)
				break

			}

			log.Debugln("Read: ", string(buf[:n]))
			out += string(buf[:n])
		}

		return out, err
	}

	return out, nil
}
