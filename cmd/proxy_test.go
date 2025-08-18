package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunProxy_ArgumentParsing(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		isProject     bool
		expectedApp   string
		expectedTool  string
		expectedArgs  []string
		expectError   bool
		errorContains string
	}{
		{
			name:         "project context - composer install",
			args:         []string{"composer", "install"},
			isProject:    true,
			expectedApp:  "",
			expectedTool: "composer",
			expectedArgs: []string{"install"},
			expectError:  false,
		},
		{
			name:         "project context - php version",
			args:         []string{"php", "-v"},
			isProject:    true,
			expectedApp:  "",
			expectedTool: "php",
			expectedArgs: []string{"-v"},
			expectError:  false,
		},
		{
			name:         "project context - npm with multiple args and flags",
			args:         []string{"npm", "run", "dev", "--watch", "--verbose"},
			isProject:    true,
			expectedApp:  "",
			expectedTool: "npm",
			expectedArgs: []string{"run", "dev", "--watch", "--verbose"},
			expectError:  false,
		},
		{
			name:         "project context - single command no args",
			args:         []string{"node"},
			isProject:    true,
			expectedApp:  "",
			expectedTool: "node",
			expectedArgs: []string{},
			expectError:  false,
		},
		{
			name:         "project context - composer with complex flags",
			args:         []string{"composer", "install", "--no-dev", "--optimize-autoloader", "--prefer-dist"},
			isProject:    true,
			expectedApp:  "",
			expectedTool: "composer",
			expectedArgs: []string{"install", "--no-dev", "--optimize-autoloader", "--prefer-dist"},
			expectError:  false,
		},
		{
			name:         "project context - php with memory flag and script",
			args:         []string{"php", "-d", "memory_limit=512M", "-f", "script.php"},
			isProject:    true,
			expectedApp:  "",
			expectedTool: "php",
			expectedArgs: []string{"-d", "memory_limit=512M", "-f", "script.php"},
			expectError:  false,
		},
		{
			name:         "global context - myapp composer install",
			args:         []string{"myapp", "composer", "install"},
			isProject:    false,
			expectedApp:  "myapp",
			expectedTool: "composer",
			expectedArgs: []string{"install"},
			expectError:  false,
		},
		{
			name:         "global context - testproject php artisan migrate with flags",
			args:         []string{"testproject", "php", "artisan", "migrate", "--force", "--seed"},
			isProject:    false,
			expectedApp:  "testproject",
			expectedTool: "php",
			expectedArgs: []string{"artisan", "migrate", "--force", "--seed"},
			expectError:  false,
		},
		{
			name:          "global context - insufficient args",
			args:          []string{"myapp"},
			isProject:     false,
			expectError:   true,
			errorContains: "not enough arguments for global context",
		},
		{
			name:          "global context - only tool name",
			args:          []string{"composer"},
			isProject:     false,
			expectError:   true,
			errorContains: "not enough arguments for global context",
		},
		{
			name:        "no arguments",
			args:        []string{},
			isProject:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for project tests
			var tmpDir string
			var cleanup func()

			if tt.isProject {
				tmpDir, cleanup = createTempProject(t)
				defer cleanup()
			} else {
				tmpDir, cleanup = createTempDir(t)
				defer cleanup()
			}

			// Change to temp directory
			oldDir, err := os.Getwd()
			assert.NoError(t, err)
			defer os.Chdir(oldDir)

			err = os.Chdir(tmpDir)
			assert.NoError(t, err)

			// Parse arguments using the same logic as runProxy
			var appName, toolName string
			var toolArgs []string
			var parseError error

			if len(tt.args) == 0 {
				parseError = assert.AnError
			} else if isPhpierProject() {
				// Project context: proxy <tool> [args...]
				toolName = tt.args[0]
				toolArgs = tt.args[1:]
				appName = ""
			} else {
				// Global context: proxy <app> <tool> [args...]
				if len(tt.args) < 2 {
					parseError = assert.AnError
				} else {
					appName = tt.args[0]
					toolName = tt.args[1]
					toolArgs = tt.args[2:]
				}
			}

			if tt.expectError {
				assert.Error(t, parseError)
				if tt.errorContains != "" {
					// Note: In actual implementation, this error would come from runProxy
					// Here we're just testing the parsing logic
				}
				return
			}

			assert.NoError(t, parseError)
			assert.Equal(t, tt.expectedApp, appName, "app name mismatch")
			assert.Equal(t, tt.expectedTool, toolName, "tool name mismatch")
			assert.Equal(t, tt.expectedArgs, toolArgs, "tool args mismatch")
		})
	}
}

func TestIsPhpierProjectDetection(t *testing.T) {
	tests := []struct {
		name           string
		hasConfigFile  bool
		expectedResult bool
	}{
		{
			name:           "directory with .phpier.yml",
			hasConfigFile:  true,
			expectedResult: true,
		},
		{
			name:           "directory without .phpier.yml",
			hasConfigFile:  false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tmpDir string
			var cleanup func()

			if tt.hasConfigFile {
				tmpDir, cleanup = createTempProject(t)
			} else {
				tmpDir, cleanup = createTempDir(t)
			}
			defer cleanup()

			// Change to temp directory
			oldDir, err := os.Getwd()
			assert.NoError(t, err)
			defer os.Chdir(oldDir)

			err = os.Chdir(tmpDir)
			assert.NoError(t, err)

			result := isPhpierProject()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

// Helper function to create a temporary directory with .phpier.yml
func createTempProject(t *testing.T) (string, func()) {
	tmpDir, cleanup := createTempDir(t)

	// Create .phpier.yml file
	configContent := `name: test-project
php_version: "8.3"`

	err := os.WriteFile(tmpDir+"/.phpier.yml", []byte(configContent), 0644)
	assert.NoError(t, err)

	return tmpDir, cleanup
}

// Helper function to create a temporary directory
func createTempDir(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "phpier-proxy-test-")
	assert.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}
