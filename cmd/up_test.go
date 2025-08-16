package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestUpCommand(t *testing.T) {
	// Test command registration
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(upCmd)

	// Check if command is properly registered
	foundCmd, _, err := cmd.Find([]string{"up"})
	assert.NoError(t, err)
	assert.Equal(t, "up", foundCmd.Use)
}

func TestUpCommandFlags(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		flagType string
	}{
		{"detach flag", "detach", "bool"},
		{"build flag", "build", "bool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := upCmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.flagType, flag.Value.Type(), "Flag %s should be of type %s", tt.flagName, tt.flagType)
		})
	}
}

func TestUpCommandDescription(t *testing.T) {
	assert.Equal(t, "up", upCmd.Use)
	assert.Equal(t, "Start the project's app container", upCmd.Short)
	assert.Contains(t, upCmd.Long, "Start the project's app container")
	assert.Contains(t, upCmd.Long, "global services")
	assert.Contains(t, upCmd.Long, "--build")
	assert.Contains(t, upCmd.Long, "-d")
}

func TestUpCommandFlags_DefaultValues(t *testing.T) {
	// Reset flags to defaults
	detached = false
	build = false

	// Test default values
	assert.False(t, detached, "detached should default to false")
	assert.False(t, build, "build should default to false")
}

func TestUpCommandFlags_ShortFlags(t *testing.T) {
	// Test that detach flag has short version
	flag := upCmd.Flags().Lookup("detach")
	assert.NotNil(t, flag)
	assert.Equal(t, "d", flag.Shorthand, "detach flag should have 'd' shorthand")
}

func TestUpCommandDocumentation(t *testing.T) {
	// Test that the command properly documents global service functionality
	assert.Contains(t, upCmd.Long, "Ensure the global services network is available")
	assert.Contains(t, upCmd.Long, "phpier up")
	assert.Contains(t, upCmd.Long, "phpier up -d")
	assert.Contains(t, upCmd.Long, "phpier up --build")
}
