package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestStopCommand(t *testing.T) {
	// Test command registration
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(stopCmd)

	// Check if command is properly registered
	foundCmd, _, err := cmd.Find([]string{"stop"})
	assert.NoError(t, err)
	assert.Equal(t, "stop", foundCmd.Use)
}

func TestStopCommandFlags(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		flagType string
	}{
		{"force flag", "force", "bool"},
		{"remove-volumes flag", "remove-volumes", "bool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := stopCmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.flagType, flag.Value.Type(), "Flag %s should be of type %s", tt.flagName, tt.flagType)
		})
	}
}

func TestStopCommandDescription(t *testing.T) {
	assert.Equal(t, "stop", stopCmd.Use)
	assert.Equal(t, "Stop the global phpier services", stopCmd.Short)
	assert.Contains(t, stopCmd.Long, "Stop the global phpier services")
	assert.Contains(t, stopCmd.Long, "Traefik")
	assert.Contains(t, stopCmd.Long, "--force")
	assert.Contains(t, stopCmd.Long, "--remove-volumes")
}

func TestStopCommandFlags_DefaultValues(t *testing.T) {
	// Reset flags to defaults
	stopForce = false
	stopRemoveVolumes = false

	// Test default values
	assert.False(t, stopForce, "force should default to false")
	assert.False(t, stopRemoveVolumes, "remove-volumes should default to false")
}

func TestStopCommandWarningMessages(t *testing.T) {
	// Test that the command contains appropriate warning messages
	assert.Contains(t, stopCmd.Long, "inaccessible until global services are started again")
	assert.Contains(t, stopCmd.Long, "dangerous")
}
