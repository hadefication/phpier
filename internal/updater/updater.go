package updater

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// AtomFeed represents the GitHub releases Atom feed
type AtomFeed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

// Entry represents a single release entry in the Atom feed
type Entry struct {
	ID      string `xml:"id"`
	Title   string `xml:"title"`
	Updated string `xml:"updated"`
}

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
	// Check for updates first
	latestVersion, err := u.getLatestVersion()
	if err != nil {
		if u.verbose {
			logrus.Debugf("Failed to check for updates: %v", err)
		}
		fmt.Println("Unable to check for updates, proceeding with install script...")
		return u.performUpdate()
	}

	// Compare versions
	if !u.needsUpdate(u.currentVer, latestVersion) {
		fmt.Printf("Already running the latest version %s\n", u.currentVer)
		return nil
	}

	fmt.Printf("Current version: %s\n", u.currentVer)
	fmt.Printf("Latest version: %s\n", latestVersion)
	fmt.Println("Updating to latest version...")

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

// getLatestVersion fetches the latest version from GitHub releases Atom feed
func (u *Updater) getLatestVersion() (string, error) {
	atomURL := "https://github.com/hadefication/phpier/releases.atom"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(atomURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch releases feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("releases feed returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read releases feed: %w", err)
	}

	var feed AtomFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return "", fmt.Errorf("failed to parse releases feed: %w", err)
	}

	if len(feed.Entries) == 0 {
		return "", fmt.Errorf("no releases found in feed")
	}

	// Extract version from the latest release title
	latestTitle := feed.Entries[0].Title
	version := u.extractVersionFromTitle(latestTitle)

	if version == "" {
		return "", fmt.Errorf("could not extract version from title: %s", latestTitle)
	}

	if u.verbose {
		logrus.Debugf("Latest version from Atom feed: %s", version)
	}

	return version, nil
}

// extractVersionFromTitle extracts version number from release title
func (u *Updater) extractVersionFromTitle(title string) string {
	// Match version patterns like "v1.2.3", "1.2.3", "Release v1.2.3", etc.
	versionRegex := regexp.MustCompile(`v?(\d+\.\d+\.\d+(?:-\w+)?)`)
	matches := versionRegex.FindStringSubmatch(title)

	if len(matches) >= 2 {
		version := matches[1]
		// Always return with 'v' prefix for consistency
		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}
		return version
	}

	return ""
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

	// Simple string comparison - in production, use proper semver library
	return currentNorm != latestNorm
}
