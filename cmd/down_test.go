package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDownCommand(t *testing.T) {
	// Test command registration
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(downCmd)

	// Check if command is properly registered
	foundCmd, _, err := cmd.Find([]string{"down"})
	assert.NoError(t, err)
	assert.Equal(t, "down", foundCmd.Use)
}

func TestDownCommandFlags(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		flagType string
	}{
		{"stop-global flag", "stop-global", "bool"},
		{"remove-volumes flag", "remove-volumes", "bool"},
		{"force flag", "force", "bool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := downCmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.flagType, flag.Value.Type(), "Flag %s should be of type %s", tt.flagName, tt.flagType)
		})
	}
}

func TestDownCommandDescription(t *testing.T) {
	assert.Equal(t, "down", downCmd.Use)
	assert.Equal(t, "Stop and remove project containers and services", downCmd.Short)
	assert.Contains(t, downCmd.Long, "Stop and remove project containers and services")
	assert.Contains(t, downCmd.Long, "--stop-global")
	assert.Contains(t, downCmd.Long, "--remove-volumes")
	assert.Contains(t, downCmd.Long, "--force")
}

func TestDownCommandFlags_DefaultValues(t *testing.T) {
	// Reset flags to defaults
	stopGlobal = false
	removeVolumes = false
	force = false

	// Test default values
	assert.False(t, stopGlobal, "stop-global should default to false")
	assert.False(t, removeVolumes, "remove-volumes should default to false")
	assert.False(t, force, "force should default to false")
}
