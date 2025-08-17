package cmd

import (
	"context"
	"os"

	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// npmCmd represents the npm command
var npmCmd = &cobra.Command{
	Use:                "npm [args...]",
	Short:              "Execute NPM commands in the app container",
	Long:               `Execute NPM commands in the app container. All arguments are forwarded to NPM.`,
	DisableFlagParsing: true,
	RunE:               runNpm,
	Example: `  phpier npm install              # Install dependencies
  phpier npm install package      # Install a specific package
  phpier npm run dev               # Run development script
  phpier npm run build             # Run build script
  phpier npm list                  # List installed packages`,
}

func runNpm(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	// Create proxy command
	proxyCmd := &docker.ProxyCommand{
		Name:        "npm",
		Command:     "npm",
		Description: "NPM package manager",
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
	rootCmd.AddCommand(npmCmd)
}
