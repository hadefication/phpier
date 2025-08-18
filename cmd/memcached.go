package cmd

import (
	"context"
	"os"

	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// memcachedCmd represents the memcached command
var memcachedCmd = &cobra.Command{
	Use:                "memcached [args...]",
	Short:              "Connect to Memcached container via telnet",
	Long:               `Connect to Memcached container via telnet interface. Provides interactive access to Memcached commands.`,
	DisableFlagParsing: true,
	RunE:               runMemcached,
	Example: `  phpier memcached                 # Launch interactive Memcached telnet session
  phpier memcached                 # Then use commands like:
                                   #   stats
                                   #   get mykey
                                   #   set mykey 0 0 5
                                   #   value
                                   #   quit`,
}

func runMemcached(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	// Create proxy command - use telnet to connect to memcached
	// Default memcached port is 11211
	telnetArgs := []string{"memcached", "11211"}
	if len(args) > 0 {
		// If user provides additional args, append them
		telnetArgs = append(telnetArgs, args...)
	}

	proxyCmd := &docker.ProxyCommand{
		Name:        "memcached",
		Command:     "telnet",
		Description: "Memcached telnet interface",
		Args:        telnetArgs,
		Interactive: true, // Telnet needs interactive TTY
	}

	// Execute the command in the global Memcached service
	exitCode, err := dockerClient.ExecuteGlobalServiceCommand(ctx, "memcached", proxyCmd)
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
	rootCmd.AddCommand(memcachedCmd)
}
