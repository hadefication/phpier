package cmd

import (
	"context"
	"os"

	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// artisanCmd represents the artisan command
var artisanCmd = &cobra.Command{
	Use:                "artisan [args...]",
	Short:              "Execute Laravel Artisan commands in the app container",
	Long:               `Execute Laravel Artisan commands in the app container. All arguments are forwarded to Artisan.`,
	DisableFlagParsing: true,
	RunE:               runArtisan,
	Example: `  phpier artisan list              # List available commands
  phpier artisan make:controller   # Create a new controller
  phpier artisan migrate           # Run database migrations
  phpier artisan tinker            # Start Artisan REPL
  phpier artisan serve             # Start development server`,
}

func runArtisan(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	// Create proxy command - Laravel artisan is executed via PHP
	proxyCmd := &docker.ProxyCommand{
		Name:        "artisan",
		Command:     "php",
		Description: "Laravel Artisan CLI",
		Args:        append([]string{"artisan"}, args...),
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
	rootCmd.AddCommand(artisanCmd)
}
