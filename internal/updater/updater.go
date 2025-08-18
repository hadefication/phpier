package updater

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// Updater handles the self-update process
type Updater struct {
	currentVer string
	verbose    bool
}

// NewUpdater creates a new updater instance
func NewUpdater(currentVersion string, verbose bool) *Updater {
	return &Updater{
		currentVer: currentVersion,
		verbose:    verbose,
	}
}

// Update performs the actual update process
func (u *Updater) Update() error {
	// Just run the update - let the install script handle everything
	return u.performUpdate()
}

// performUpdate executes the update using the install script
func (u *Updater) performUpdate() error {
	// Determine current install directory
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	installDir := strings.TrimSuffix(execPath, "/phpier")
	if u.verbose {
		logrus.Debugf("Current installation directory: %s", installDir)
	}

	// Build install script command
	installScriptURL := "https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh"

	// Prepare arguments for the install script
	args := []string{"-s", "--", "--dir", installDir, "--force"}

	if u.verbose {
		logrus.Debugf("Running install script with args: %v", args)
	}

	// Execute install script via curl and bash
	curlCmd := fmt.Sprintf("curl -sSL %s | bash %s", installScriptURL, strings.Join(args, " "))

	cmd := exec.Command("bash", "-c", curlCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment variables for the install script
	cmd.Env = append(os.Environ(),
		"FORCE_INSTALL=true", // Skip confirmation prompts
	)

	// Run the installation - let install script handle version comparison and messaging
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}
