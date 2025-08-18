package cmd

import (
	"context"
	"os"

	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// redisCmd represents the redis command
var redisCmd = &cobra.Command{
	Use:                "redis [args...]",
	Short:              "Execute Redis CLI commands in the Redis/Valkey container",
	Long:               `Execute Redis CLI commands in the Redis/Valkey container. All arguments are forwarded to redis-cli.`,
	DisableFlagParsing: true,
	RunE:               runRedis,
	Example: `  phpier redis                     # Launch interactive Redis CLI
  phpier redis ping                # Test Redis connection
  phpier redis keys "*"            # List all keys
  phpier redis get mykey           # Get value of a key
  phpier redis set mykey value     # Set a key-value pair`,
}

func runRedis(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	// Create proxy command
	proxyCmd := &docker.ProxyCommand{
		Name:        "redis",
		Command:     "redis-cli",
		Description: "Redis CLI",
		Args:        args,
		Interactive: len(args) == 0, // Only interactive if no args provided (entering Redis shell)
	}

	// Execute the command in the global Redis service
	exitCode, err := dockerClient.ExecuteGlobalServiceCommand(ctx, "redis", proxyCmd)
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
	rootCmd.AddCommand(redisCmd)
}
