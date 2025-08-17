package cmd

import (
	"context"
	"os"

	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// phpCmd represents the php command
var phpCmd = &cobra.Command{
	Use:                "php [args...]",
	Short:              "Execute PHP commands in the app container",
	Long:               `Execute PHP commands in the app container. All arguments are forwarded to the PHP interpreter.`,
	DisableFlagParsing: true,
	RunE:               runPhp,
	Example: `  phpier php -v                    # Show PHP version
  phpier php -m                    # List loaded modules
  phpier php script.php            # Run a PHP script
  phpier php -r "echo 'Hello';"    # Execute PHP code directly`,
}

func runPhp(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	// Create proxy command
	proxyCmd := &docker.ProxyCommand{
		Name:        "php",
		Command:     "php",
		Description: "PHP interpreter",
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
	rootCmd.AddCommand(phpCmd)
}
