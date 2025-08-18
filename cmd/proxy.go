package cmd

import (
	"context"
	"fmt"
	"os"

	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// proxyCmd represents the unified proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy [app] <tool> [args...]",
	Short: "Execute tools in app containers with context-aware behavior",
	Long: `Execute tools in app containers with context-aware behavior.
All flags and arguments are automatically forwarded to the target tool.

Context-aware usage:
• In phpier project: phpier proxy <tool> [args...]    - executes in current project's app container
• Outside project:   phpier proxy <app> <tool> [args...] - executes in specified app's container

Examples:
  phpier proxy composer install --no-dev           # Install dependencies with flags
  phpier proxy php -v                              # Show PHP version
  phpier proxy npm run dev -- --watch              # Run npm script with arguments
  phpier proxy php -d memory_limit=512M script.php # PHP with configuration flags
  phpier proxy myapp composer require --dev phpunit/phpunit  # Global context with flags
  phpier proxy myapp php artisan migrate --force   # Laravel migration with force flag`,
	DisableFlagParsing: true,
	RunE:               runProxy,
	Args:               cobra.MinimumNArgs(1),
}

func runProxy(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	var appName, toolName string
	var toolArgs []string

	// Detect context and parse arguments
	// Note: DisableFlagParsing: true ensures all flags are passed through unchanged
	if isPhpierProject() {
		// Project context: proxy <tool> [args...]
		// All arguments after tool name (including flags) are forwarded
		toolName = args[0]
		toolArgs = args[1:]
		appName = "" // Will be determined from project config
	} else {
		// Global context: proxy <app> <tool> [args...]
		// All arguments after tool name (including flags) are forwarded
		if len(args) < 2 {
			return fmt.Errorf(`not enough arguments for global context

Usage:
• In phpier project: phpier proxy <tool> [args...]
• Outside project:   phpier proxy <app> <tool> [args...]

Examples:
  phpier proxy composer install           # In project directory
  phpier proxy myapp composer install     # From anywhere`)
		}
		appName = args[0]
		toolName = args[1]
		toolArgs = args[2:]
	}

	// Create proxy command
	proxyCommand := &docker.ProxyCommand{
		Name:        toolName,
		Command:     toolName,
		Description: fmt.Sprintf("%s command", toolName),
		Args:        toolArgs,
	}

	// Execute based on context
	var exitCode int
	if appName == "" {
		// Project context - use existing ExecuteProxyCommand
		exitCode, err = dockerClient.ExecuteProxyCommand(ctx, proxyCommand)
	} else {
		// Global context - execute in specified app container
		exitCode, err = dockerClient.ExecuteGlobalProxyCommand(ctx, appName, proxyCommand)
	}

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
	rootCmd.AddCommand(proxyCmd)
}
