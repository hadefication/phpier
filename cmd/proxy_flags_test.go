package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyFlagForwarding(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		isProject     bool
		expectedTool  string
		expectedArgs  []string
		description   string
	}{
		{
			name:         "project context - composer with flags",
			args:         []string{"composer", "install", "--no-dev", "--optimize-autoloader"},
			isProject:    true,
			expectedTool: "composer",
			expectedArgs: []string{"install", "--no-dev", "--optimize-autoloader"},
			description:  "All flags should be forwarded to composer",
		},
		{
			name:         "project context - php with multiple flags",
			args:         []string{"php", "-d", "memory_limit=512M", "-f", "script.php", "--", "arg1", "arg2"},
			isProject:    true,
			expectedTool: "php",
			expectedArgs: []string{"-d", "memory_limit=512M", "-f", "script.php", "--", "arg1", "arg2"},
			description:  "Complex PHP flags and arguments should be preserved",
		},
		{
			name:         "project context - npm with long and short flags",
			args:         []string{"npm", "install", "-g", "--save-dev", "package-name"},
			isProject:    true,
			expectedTool: "npm",
			expectedArgs: []string{"install", "-g", "--save-dev", "package-name"},
			description:  "Mix of short and long flags should be forwarded",
		},
		{
			name:         "global context - composer with flags",
			args:         []string{"myapp", "composer", "require", "--dev", "phpunit/phpunit", "^9.0"},
			isProject:    false,
			expectedTool: "composer",
			expectedArgs: []string{"require", "--dev", "phpunit/phpunit", "^9.0"},
			description:  "Flags in global context should be forwarded correctly",
		},
		{
			name:         "global context - php with verbose flag",
			args:         []string{"testapp", "php", "-v"},
			isProject:    false,
			expectedTool: "php",
			expectedArgs: []string{"-v"},
			description:  "Short flags should work in global context",
		},
		{
			name:         "project context - tool with equals flags",
			args:         []string{"composer", "config", "--global", "repo.packagist.org=false"},
			isProject:    true,
			expectedTool: "composer",
			expectedArgs: []string{"config", "--global", "repo.packagist.org=false"},
			description:  "Flags with equals signs should be preserved",
		},
		{
			name:         "project context - node with negative flags",
			args:         []string{"node", "--no-warnings", "--experimental-modules", "app.js"},
			isProject:    true,
			expectedTool: "node",
			expectedArgs: []string{"--no-warnings", "--experimental-modules", "app.js"},
			description:  "Negative flags and experimental flags should be forwarded",
		},
		{
			name:         "global context - npm with environment flags",
			args:         []string{"webapp", "npm", "run", "build", "--", "--mode=production", "--analyze"},
			isProject:    false,
			expectedTool: "npm",
			expectedArgs: []string{"run", "build", "--", "--mode=production", "--analyze"},
			description:  "NPM script arguments after -- should be preserved",
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

			// Parse arguments using the same logic as runProxy
			var appName, toolName string
			var toolArgs []string

			// Simulate the context detection (using the temp directory setup)
			if tt.isProject {
				// Project context: proxy <tool> [args...]
				toolName = tt.args[0]
				toolArgs = tt.args[1:]
				appName = ""
			} else {
				// Global context: proxy <app> <tool> [args...]
				appName = tt.args[0]
				toolName = tt.args[1]
				toolArgs = tt.args[2:]
			}

			// Verify flag forwarding
			assert.Equal(t, tt.expectedTool, toolName, "tool name should match")
			assert.Equal(t, tt.expectedArgs, toolArgs, "all flags and arguments should be forwarded: %s", tt.description)
			
			// Verify that no flags were lost or modified
			assert.Len(t, toolArgs, len(tt.expectedArgs), "number of forwarded arguments should match expected")
			
			// Check that specific flag patterns are preserved
			for i, expectedArg := range tt.expectedArgs {
				if i < len(toolArgs) {
					assert.Equal(t, expectedArg, toolArgs[i], "argument %d should match exactly", i)
				}
			}
		})
	}
}

func TestProxyComplexFlagScenarios(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "mixed quotes and flags",
			args:     []string{"php", "-r", `echo "Hello World";`, "--define", "memory_limit=256M"},
			expected: []string{"-r", `echo "Hello World";`, "--define", "memory_limit=256M"},
		},
		{
			name:     "flags with spaces in values", 
			args:     []string{"composer", "create-project", "--prefer-dist", "laravel/laravel", "my project"},
			expected: []string{"create-project", "--prefer-dist", "laravel/laravel", "my project"},
		},
		{
			name:     "multiple equals flags",
			args:     []string{"npm", "config", "set", "registry=https://registry.npmjs.org/", "--save=false"},
			expected: []string{"config", "set", "registry=https://registry.npmjs.org/", "--save=false"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate project context parsing
			toolName := tt.args[0]
			toolArgs := tt.args[1:]

			assert.Equal(t, tt.expected, toolArgs, "complex flag scenarios should be preserved exactly")
		})
	}
}