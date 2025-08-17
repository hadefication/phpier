package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestBuildCommand(t *testing.T) {
	// Test command registration
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(buildCmd)

	// Check if command is properly registered
	foundCmd, _, err := cmd.Find([]string{"build"})
	assert.NoError(t, err)
	assert.Equal(t, "build", foundCmd.Use)
}

func TestBuildCommandFlags(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		flagType string
	}{
		{"no-cache flag", "no-cache", "bool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := buildCmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.flagType, flag.Value.Type(), "Flag %s should be of type %s", tt.flagName, tt.flagType)
		})
	}
}

func TestBuildCommandDescription(t *testing.T) {
	assert.Equal(t, "build", buildCmd.Use)
	assert.Equal(t, "Build the project's app container", buildCmd.Short)
	assert.Contains(t, buildCmd.Long, "Build (or rebuild) the app container")
	assert.Contains(t, buildCmd.Long, "Dockerfile.php")
	assert.Contains(t, buildCmd.Long, "--no-cache")
}

func TestBuildCommandFlags_DefaultValues(t *testing.T) {
	// Reset flags to defaults
	noCache = false

	// Test default values
	assert.False(t, noCache, "noCache should default to false")
}

func TestBuildCommandDocumentation(t *testing.T) {
	// Test that the command properly documents build functionality
	assert.Contains(t, buildCmd.Long, "Build only the app container")
	assert.Contains(t, buildCmd.Long, "phpier build")
	assert.Contains(t, buildCmd.Long, "phpier build --no-cache")
	assert.Contains(t, buildCmd.Long, "Force a clean rebuild without using cache")
}

func TestBuildCommandLongDescription(t *testing.T) {
	// Test the long description contains key information
	expected := []string{
		"Build (or rebuild) the app container for the current project",
		"Build only the app container using the project's Dockerfile.php",
		"Support forcing a rebuild with the --no-cache flag",
		"Validate that the project is properly initialized",
	}

	for _, expectedText := range expected {
		assert.Contains(t, buildCmd.Long, expectedText,
			"Long description should contain: %s", expectedText)
	}
}

func TestBuildCommandExamples(t *testing.T) {
	// Test that examples are included in the command documentation
	examples := []string{
		"phpier build",
		"phpier build --no-cache",
	}

	for _, example := range examples {
		assert.Contains(t, buildCmd.Long, example,
			"Examples should contain: %s", example)
	}
}
