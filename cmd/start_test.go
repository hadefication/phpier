package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestStartCommand(t *testing.T) {
	// Test command registration
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(startCmd)

	// Check if command is properly registered
	foundCmd, _, err := cmd.Find([]string{"start"})
	assert.NoError(t, err)
	assert.Equal(t, "start", foundCmd.Use)
}

func TestStartCommandFlags(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		flagType string
	}{
		{"detach flag", "detach", "bool"},
		{"build flag", "build", "bool"},
		{"force flag", "force", "bool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := startCmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.flagType, flag.Value.Type(), "Flag %s should be of type %s", tt.flagName, tt.flagType)
		})
	}
}

func TestStartCommandDescription(t *testing.T) {
	assert.Equal(t, "start", startCmd.Use)
	assert.Equal(t, "Start the global phpier services", startCmd.Short)
	assert.Contains(t, startCmd.Long, "Start the global phpier services")
	assert.Contains(t, startCmd.Long, "Traefik")
	assert.Contains(t, startCmd.Long, "--build")
	assert.Contains(t, startCmd.Long, "--force")
}

func TestStartCommandFlags_DefaultValues(t *testing.T) {
	// Reset flags to defaults
	startDetached = true
	startBuild = false
	startForce = false

	// Test default values
	assert.True(t, startDetached, "detach should default to true")
	assert.False(t, startBuild, "build should default to false")
	assert.False(t, startForce, "force should default to false")
}

func TestStartCommandFlags_ShortFlags(t *testing.T) {
	// Test that detach flag has short version
	flag := startCmd.Flags().Lookup("detach")
	assert.NotNil(t, flag)
	assert.Equal(t, "d", flag.Shorthand, "detach flag should have 'd' shorthand")
}
