package cmd

import (
	"context"
	"os"

	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// npxCmd represents the npx command
var npxCmd = &cobra.Command{
	Use:                "npx [args...]",
	Short:              "Execute NPX commands in the app container",
	Long:               `Execute NPX commands in the app container. All arguments are forwarded to NPX.`,
	DisableFlagParsing: true,
	RunE:               runNpx,
	Example: `  phpier npx create-react-app myapp  # Create a new React app
  phpier npx vite                    # Run Vite
  phpier npx eslint .                # Run ESLint
  phpier npx prettier --write .      # Format code with Prettier
  phpier npx --version               # Show NPX version`,
}

func runNpx(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	// Create proxy command
	proxyCmd := &docker.ProxyCommand{
		Name:        "npx",
		Command:     "npx",
		Description: "NPX package runner",
		Args:        args,
	}

	// Execute the command
	exitCode, err := dockerClient.ExecuteProxyCommand(ctx, proxyCmd)
	if err != nil {
		return err
	}

	// Exit with the same code as the container command
	if exitCode != 0 {
		os.Exit(exitCode)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(npxCmd)
}
