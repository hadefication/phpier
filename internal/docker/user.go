package docker

import (
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// SetWWWUser sets WWWUSER environment variable to current user ID if not already set.
// This ensures file permissions work correctly in Docker containers.
func SetWWWUser() error {
	if os.Getenv("WWWUSER") != "" {
		return nil // Already set
	}

	// Try using $UID first, fallback to id -u command
	uid := os.Getenv("UID")
	if uid == "" {
		cmd := exec.Command("id", "-u")
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		uid = strings.TrimSpace(string(output))
	}

	os.Setenv("WWWUSER", uid)
	logrus.Debugf("Set WWWUSER=%s (current user ID)", uid)
	return nil
}
