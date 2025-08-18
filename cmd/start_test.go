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
	// The start command should have no flags (simple implementation)
	flags := startCmd.Flags()
	assert.Equal(t, 0, flags.NFlag(), "Start command should have no flags")
}

func TestStartCommandDescription(t *testing.T) {
	assert.Equal(t, "start", startCmd.Use)
	assert.Equal(t, "Start services (global only if not in project, global + project if in project)", startCmd.Short)
	assert.Contains(t, startCmd.Long, "Start services based on current directory context")
	assert.Contains(t, startCmd.Long, "When run outside a phpier project directory")
	assert.Contains(t, startCmd.Long, "When run inside a phpier project directory")
	assert.Contains(t, startCmd.Long, "Context-aware service startup")
}
