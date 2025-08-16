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
	// The stop command should have no flags (simple implementation)
	flags := stopCmd.Flags()
	assert.Equal(t, 0, flags.NFlag(), "Stop command should have no flags")
}

func TestStopCommandDescription(t *testing.T) {
	assert.Equal(t, "stop", stopCmd.Use)
	assert.Equal(t, "Stop project containers and global services (shortcut for 'down --global')", stopCmd.Short)
	assert.Contains(t, stopCmd.Long, "Stop and remove project containers and services, then stop global services")
	assert.Contains(t, stopCmd.Long, "shortcut for 'phpier down --global'")
	assert.Contains(t, stopCmd.Long, "phpier stop")
}
