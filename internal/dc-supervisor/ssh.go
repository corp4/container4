package supervisor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// Structure that holds RPC methods
type SSH struct{}

// Adds a public key to the authorized_keys file of the current user
// The pubKey is added as is, without any validation
func (s *SSH) AddAuthorizedKey(pubKey string, success *bool) error {
	// Check if public key is already in authorized_keys file
	if s.HasAuthorizedKey(pubKey, success); *success {
		return nil
	}

	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Check if .ssh directory exists
	if _, err := os.Stat(filepath.Join(homeDir, ".ssh")); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join(homeDir, ".ssh"), 0700); err != nil {
			return fmt.Errorf("failed to create .ssh directory: %w", err)
		}
	}

	// Open authorized_keys file
	file, err := os.OpenFile(filepath.Join(homeDir, ".ssh", "authorized_keys"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open authorized_keys file: %w", err)
	}
	defer file.Close()

	// Write public key to authorized_keys file
	if _, err := file.WriteString(pubKey + "\n"); err != nil {
		return fmt.Errorf("failed to write public key to authorized_keys file: %w", err)
	}

	*success = true
	return nil
}

// Checks if a public key is in the authorized_keys file of the current user
func (s *SSH) HasAuthorizedKey(pubkey string, hasKey *bool) error {
	*hasKey = false

	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Open authorized_keys file
	file, err := os.ReadFile(filepath.Join(homeDir, ".ssh", "authorized_keys"))
	if err != nil {
		return fmt.Errorf("failed to open authorized_keys file: %w", err)
	}

	// Check if public key is in authorized_keys file
	*hasKey = bytes.Contains(file, []byte(pubkey))
	return nil
}
