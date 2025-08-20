package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	forceUninstall bool
	dryRun         bool
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall phpier from your system",
	Long: `Uninstall phpier from your system by running the uninstall script.

This command will:
- Remove the phpier binary from common installation paths
- Clean up phpier configuration files and data directories
- Stop and remove phpier Docker containers and networks
- Remove built phpier project images (phpier-* images)
- Preserve Docker volumes for data safety

Examples:
  phpier uninstall          # Interactive uninstall with confirmation
  phpier uninstall --force  # Skip confirmation prompts
  phpier uninstall --dry-run # Show what would be removed without removing`,
	RunE: runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Flags
	uninstallCmd.Flags().BoolVarP(&forceUninstall, "force", "f", false, "Skip confirmation prompts")
	uninstallCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be removed without actually removing")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	logrus.Infof("🗑️  Starting phpier uninstallation...")

	// Find the uninstall script
	scriptPath, err := findUninstallScript()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeFileNotFound, "Failed to find uninstall script", err)
	}

	logrus.Debugf("Found uninstall script at: %s", scriptPath)

	// Build arguments for the uninstall script
	args = []string{}
	if forceUninstall {
		args = append(args, "--force")
	}
	if dryRun {
		args = append(args, "--dry-run")
	}

	logrus.Infof("🚀 Running uninstall script...")
	logrus.Debugf("Executing: %s %v", scriptPath, args)

	// Execute the uninstall script
	execCmd := exec.Command(scriptPath, args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	if err := execCmd.Run(); err != nil {
		return errors.WrapError(errors.ErrorTypeCommandFailed, "Uninstall script failed", err)
	}

	return nil
}

// findUninstallScript locates the uninstall.sh script
func findUninstallScript() (string, error) {
	// Get the directory where the current phpier binary is located
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	execDir := filepath.Dir(execPath)

	// Try multiple possible locations for the uninstall script
	possiblePaths := []string{
		// If running from development environment
		filepath.Join(execDir, "scripts", "uninstall.sh"),
		// If installed and scripts are in same directory
		filepath.Join(execDir, "uninstall.sh"),
		// If scripts are in a sibling directory
		filepath.Join(filepath.Dir(execDir), "scripts", "uninstall.sh"),
		// If running from source directory
		filepath.Join(execDir, "..", "scripts", "uninstall.sh"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			// Check if it's executable
			if info, err := os.Stat(path); err == nil && info.Mode()&0111 != 0 {
				return path, nil
			}
		}
	}

	return "", errors.NewPhpierError(errors.ErrorTypeFileNotFound, "Uninstall script not found").
		WithSuggestion("Make sure phpier was installed correctly").
		WithSuggestion("Try downloading the uninstall script manually from the repository")
}
