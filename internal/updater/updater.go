package updater

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// UpdateOptions holds configuration for update process
type UpdateOptions struct {
	Version   string // Specific version to update to (empty for latest)
	CheckOnly bool   // Only check for updates, don't install
	Force     bool   // Skip confirmation prompts
	Verbose   bool   // Show verbose output
}

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	CurrentVersion string
	LatestVersion  string
	ReleaseNotes   string
	UpdateURL      string
	AssetSize      int64
	NeedsUpdate    bool
}

// Updater handles the self-update process
type Updater struct {
	github     *GitHubClient
	currentVer string
	verbose    bool
}

// NewUpdater creates a new updater instance
func NewUpdater(currentVersion string, verbose bool) *Updater {
	return &Updater{
		github:     NewGitHubClient(),
		currentVer: currentVersion,
		verbose:    verbose,
	}
}

// CheckForUpdates checks if an update is available
func (u *Updater) CheckForUpdates(targetVersion string) (*UpdateInfo, error) {
	var release *GitHubRelease
	var err error

	if targetVersion != "" {
		// Check for specific version
		release, err = u.github.GetRelease(targetVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to get release %s: %w", targetVersion, err)
		}
	} else {
		// Check for latest version
		release, err = u.github.GetLatestRelease()
		if err != nil {
			return nil, fmt.Errorf("failed to get latest release: %w", err)
		}
	}

	// Find asset for current platform
	asset, err := release.FindAssetForPlatform(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, fmt.Errorf("no compatible release found: %w", err)
	}

	info := &UpdateInfo{
		CurrentVersion: u.currentVer,
		LatestVersion:  release.TagName,
		ReleaseNotes:   release.Body,
		UpdateURL:      asset.BrowserDownloadURL,
		AssetSize:      asset.Size,
		NeedsUpdate:    u.needsUpdate(u.currentVer, release.TagName),
	}

	return info, nil
}

// Update performs the actual update process
func (u *Updater) Update(options UpdateOptions) error {
	// Check for updates first
	updateInfo, err := u.CheckForUpdates(options.Version)
	if err != nil {
		return err
	}

	if options.CheckOnly {
		return u.printUpdateInfo(updateInfo)
	}

	if !updateInfo.NeedsUpdate {
		fmt.Printf("Already running latest version %s\n", updateInfo.CurrentVersion)
		return nil
	}

	// Show update information
	fmt.Printf("Current version: %s\n", updateInfo.CurrentVersion)
	fmt.Printf("Latest version:  %s\n", updateInfo.LatestVersion)
	fmt.Printf("Download size:   %s\n", formatBytes(updateInfo.AssetSize))

	// Confirm update unless forced
	if !options.Force {
		fmt.Print("\nProceed with update? [y/N]: ")
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Update cancelled")
			return nil
		}
	}

	return u.performUpdate(updateInfo, options.Verbose)
}

// performUpdate executes the update using the install script
func (u *Updater) performUpdate(info *UpdateInfo, verbose bool) error {
	fmt.Printf("🔄 Updating phpier to version %s...\n", info.LatestVersion)

	// Determine current install directory
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	installDir := strings.TrimSuffix(execPath, "/phpier")
	if verbose {
		logrus.Debugf("Current installation directory: %s", installDir)
	}

	// Build install script command
	installScriptURL := "https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh"

	// Prepare arguments for the install script
	args := []string{"-s", "--"}

	if info.LatestVersion != "" {
		args = append(args, "--version", info.LatestVersion)
	}

	args = append(args, "--dir", installDir, "--force")

	if verbose {
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

	// Run the installation
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run install script: %w", err)
	}

	fmt.Printf("✅ Successfully updated to version %s\n", info.LatestVersion)
	fmt.Println("Run 'phpier version' to verify the update")

	return nil
}

// needsUpdate determines if current version needs update
func (u *Updater) needsUpdate(current, latest string) bool {
	// Normalize versions by removing 'v' prefix if present
	currentNorm := strings.TrimPrefix(current, "v")
	latestNorm := strings.TrimPrefix(latest, "v")

	// If latest version is empty or unknown, no update available
	if latestNorm == "" || latestNorm == "unknown" {
		return false
	}

	// If current version is empty or unknown, update is needed
	if currentNorm == "" || currentNorm == "unknown" {
		return true
	}

	// Simple comparison - in production, use proper semver library
	return currentNorm != latestNorm
}

// printUpdateInfo displays update information for check-only mode
func (u *Updater) printUpdateInfo(info *UpdateInfo) error {
	fmt.Printf("Current version: %s\n", info.CurrentVersion)
	fmt.Printf("Latest version:  %s\n", info.LatestVersion)

	if info.NeedsUpdate {
		fmt.Printf("📦 Update available! Download size: %s\n", formatBytes(info.AssetSize))
		fmt.Println("Run 'phpier self-update' to install the update")

		if info.ReleaseNotes != "" {
			fmt.Println("\nRelease notes:")
			fmt.Println(strings.TrimSpace(info.ReleaseNotes))
		}
	} else {
		fmt.Println("✅ You're running the latest version")
	}

	return nil
}

// IsValidVersionFormat checks if the version string is in a valid format
func IsValidVersionFormat(version string) bool {
	// Allow both v1.2.3 and 1.2.3 formats
	// Simple validation - in production, use semver library
	if len(version) == 0 {
		return false
	}

	// Must start with v or digit
	if version[0] != 'v' && (version[0] < '0' || version[0] > '9') {
		return false
	}

	// Basic format validation
	// For now, just check it's not empty and starts correctly
	return len(version) > 1
}
