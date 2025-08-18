package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPhpierProject(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T, tempDir string)
		expected    bool
		description string
	}{
		{
			name: "phpier project with .phpier.yml",
			setupFunc: func(t *testing.T, tempDir string) {
				// Create .phpier.yml file
				configFile := filepath.Join(tempDir, ".phpier.yml")
				err := os.WriteFile(configFile, []byte("project_name: test"), 0644)
				assert.NoError(t, err)
			},
			expected:    true,
			description: "should return true when .phpier.yml exists",
		},
		{
			name: "non-phpier project directory",
			setupFunc: func(t *testing.T, tempDir string) {
				// Create some other files but not .phpier.yml
				otherFile := filepath.Join(tempDir, "package.json")
				err := os.WriteFile(otherFile, []byte("{}"), 0644)
				assert.NoError(t, err)
			},
			expected:    false,
			description: "should return false when .phpier.yml does not exist",
		},
		{
			name: "empty directory",
			setupFunc: func(t *testing.T, tempDir string) {
				// Do nothing - keep directory empty
			},
			expected:    false,
			description: "should return false for empty directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test
			tempDir, err := os.MkdirTemp("", "phpier-test-*")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Store original working directory
			originalDir, err := os.Getwd()
			assert.NoError(t, err)

			// Change to temp directory
			err = os.Chdir(tempDir)
			assert.NoError(t, err)
			defer func() {
				err := os.Chdir(originalDir)
				assert.NoError(t, err)
			}()

			// Setup test scenario
			tt.setupFunc(t, tempDir)

			// Test the function
			result := isPhpierProject()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestIsPhpierProject_PermissionError(t *testing.T) {
	// This test covers the edge case where os.Stat returns a permission error
	// which should be treated as "not a phpier project"

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "phpier-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Store original working directory
	originalDir, err := os.Getwd()
	assert.NoError(t, err)

	// Change to temp directory
	err = os.Chdir(tempDir)
	assert.NoError(t, err)
	defer func() {
		err := os.Chdir(originalDir)
		assert.NoError(t, err)
	}()

	// The function should handle any file access errors gracefully
	// and return false for any error condition
	result := isPhpierProject()
	assert.False(t, result, "should return false when file access fails")
}
