package cmd

import (
	"context"
	"os"

	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:                "node [args...]",
	Short:              "Execute Node.js commands in the app container",
	Long:               `Execute Node.js commands in the app container. All arguments are forwarded to Node.js.`,
	DisableFlagParsing: true,
	RunE:               runNode,
	Example: `  phpier node -v                   # Show Node.js version
  phpier node script.js            # Run a JavaScript file
  phpier node -e "console.log('Hi')" # Execute JavaScript code directly
  phpier node --help               # Show Node.js help`,
}

func runNode(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	// Create proxy command
	proxyCmd := &docker.ProxyCommand{
		Name:        "node",
		Command:     "node",
		Description: "Node.js runtime",
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
	rootCmd.AddCommand(nodeCmd)
}
