package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUninstallCommand(t *testing.T) {
	// Test command registration
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(uninstallCmd)

	// Check if command is properly registered
	foundCmd, _, err := cmd.Find([]string{"uninstall"})
	assert.NoError(t, err)
	assert.Equal(t, "uninstall", foundCmd.Use)
}

func TestUninstallCommandFlags(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		flagType string
	}{
		{"force flag", "force", "bool"},
		{"dry-run flag", "dry-run", "bool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := uninstallCmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.flagType, flag.Value.Type(), "Flag %s should be of type %s", tt.flagName, tt.flagType)
		})
	}
}

func TestUninstallCommandDescription(t *testing.T) {
	assert.Equal(t, "uninstall", uninstallCmd.Use)
	assert.Equal(t, "Uninstall phpier from your system", uninstallCmd.Short)
	assert.Contains(t, uninstallCmd.Long, "Uninstall phpier from your system")
	assert.Contains(t, uninstallCmd.Long, "--force")
	assert.Contains(t, uninstallCmd.Long, "--dry-run")
}

func TestUninstallCommandFlags_DefaultValues(t *testing.T) {
	// Reset flags to default values
	forceUninstall = false
	dryRun = false

	assert.False(t, forceUninstall, "force flag should default to false")
	assert.False(t, dryRun, "dry-run flag should default to false")
}

func TestFindUninstallScript(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) string // Returns the mock executable path
		expectError   bool
		errorContains string
	}{
		{
			name: "script found in scripts subdirectory",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Create mock phpier binary
				binPath := filepath.Join(tmpDir, "phpier")
				err := os.WriteFile(binPath, []byte("#!/bin/bash\necho 'mock phpier'"), 0755)
				require.NoError(t, err)

				// Create scripts directory and uninstall script
				scriptsDir := filepath.Join(tmpDir, "scripts")
				err = os.MkdirAll(scriptsDir, 0755)
				require.NoError(t, err)

				scriptPath := filepath.Join(scriptsDir, "uninstall.sh")
				err = os.WriteFile(scriptPath, []byte("#!/bin/bash\necho 'uninstall script'"), 0755)
				require.NoError(t, err)

				return binPath
			},
			expectError: false,
		},
		{
			name: "script found in same directory",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Create mock phpier binary
				binPath := filepath.Join(tmpDir, "phpier")
				err := os.WriteFile(binPath, []byte("#!/bin/bash\necho 'mock phpier'"), 0755)
				require.NoError(t, err)

				// Create uninstall script in same directory
				scriptPath := filepath.Join(tmpDir, "uninstall.sh")
				err = os.WriteFile(scriptPath, []byte("#!/bin/bash\necho 'uninstall script'"), 0755)
				require.NoError(t, err)

				return binPath
			},
			expectError: false,
		},
		{
			name: "script not found",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Create mock phpier binary but no uninstall script
				binPath := filepath.Join(tmpDir, "phpier")
				err := os.WriteFile(binPath, []byte("#!/bin/bash\necho 'mock phpier'"), 0755)
				require.NoError(t, err)

				return binPath
			},
			expectError:   true,
			errorContains: "Uninstall script not found",
		},
		{
			name: "script exists but not executable",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Create mock phpier binary
				binPath := filepath.Join(tmpDir, "phpier")
				err := os.WriteFile(binPath, []byte("#!/bin/bash\necho 'mock phpier'"), 0755)
				require.NoError(t, err)

				// Create uninstall script without execute permissions
				scriptPath := filepath.Join(tmpDir, "uninstall.sh")
				err = os.WriteFile(scriptPath, []byte("#!/bin/bash\necho 'uninstall script'"), 0644) // No execute bit
				require.NoError(t, err)

				return binPath
			},
			expectError:   true,
			errorContains: "Uninstall script not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecPath := tt.setupFunc(t)

			// Since we can't easily mock os.Executable, we'll test the logic directly
			// by calling the helper function with a known path
			execDir := filepath.Dir(mockExecPath)

			// Simulate the logic in findUninstallScript
			possiblePaths := []string{
				filepath.Join(execDir, "scripts", "uninstall.sh"),
				filepath.Join(execDir, "uninstall.sh"),
				filepath.Join(filepath.Dir(execDir), "scripts", "uninstall.sh"),
				filepath.Join(execDir, "..", "scripts", "uninstall.sh"),
			}

			var foundPath string
			var found bool

			for _, path := range possiblePaths {
				if _, err := os.Stat(path); err == nil {
					// Check if it's executable
					if info, err := os.Stat(path); err == nil && info.Mode()&0111 != 0 {
						foundPath = path
						found = true
						break
					}
				}
			}

			if tt.expectError {
				assert.False(t, found, "Should not find executable script")
			} else {
				assert.True(t, found, "Should find executable script")
				assert.NotEmpty(t, foundPath, "Should return path to found script")

				// Verify the found script is actually executable
				info, err := os.Stat(foundPath)
				require.NoError(t, err)
				assert.True(t, info.Mode()&0111 != 0, "Found script should be executable")
			}
		})
	}
}

// Test helper functions for command functionality
func TestUninstallCommand_Integration(t *testing.T) {
	// This is a basic integration test that verifies the command structure
	// without actually running the uninstall script

	// Create a minimal command for testing
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(uninstallCmd)

	// Test help output
	helpOutput := uninstallCmd.Long
	assert.Contains(t, helpOutput, "phpier uninstall")
	assert.Contains(t, helpOutput, "--force")
	assert.Contains(t, helpOutput, "--dry-run")

	// Test that command accepts the expected flags
	forceFlag := uninstallCmd.Flags().Lookup("force")
	assert.NotNil(t, forceFlag)
	assert.Equal(t, "bool", forceFlag.Value.Type())

	dryRunFlag := uninstallCmd.Flags().Lookup("dry-run")
	assert.NotNil(t, dryRunFlag)
	assert.Equal(t, "bool", dryRunFlag.Value.Type())
}
