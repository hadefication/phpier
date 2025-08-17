package cmd

import (
	"context"
	"os"

	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// composerCmd represents the composer command
var composerCmd = &cobra.Command{
	Use:                "composer [args...]",
	Short:              "Execute Composer commands in the app container",
	Long:               `Execute Composer commands in the app container. All arguments are forwarded to Composer.`,
	DisableFlagParsing: true,
	RunE:               runComposer,
	Example: `  phpier composer install          # Install dependencies
  phpier composer require package  # Add a new package
  phpier composer update           # Update dependencies
  phpier composer show             # Show installed packages
  phpier composer dump-autoload    # Regenerate autoloader`,
}

func runComposer(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	// Create proxy command
	proxyCmd := &docker.ProxyCommand{
		Name:        "composer",
		Command:     "composer",
		Description: "Composer dependency manager",
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
	rootCmd.AddCommand(composerCmd)
}
