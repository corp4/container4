package supervisor

import (
	"time"
)

// Structure that holds RPC methods
type Supervisor struct{}

// Used to check if the supervisor is running
func (s *Supervisor) GetStatus(_ struct{}, status *string) error {
	*status = "UP"
	return nil
}

// Returns the current time on the devcontainer
func (s *Supervisor) GetTime(_ struct{}, t *time.Time) error {
	*t = time.Now()
	return nil
}
