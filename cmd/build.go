package cmd

import (
	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"
	"phpier/internal/generator"
	"phpier/internal/templates"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	noCache    bool
	regenerate bool
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project's app container",
	Long: `Build (or rebuild) the app container for the current project.

This command will:
- Build only the app container using the project's Dockerfile.php
- Support forcing a rebuild with the --no-cache flag
- Validate that the project is properly initialized
- Use existing configuration files (preserves customizations)

Examples:
  phpier build                    # Build the app container using existing files
  phpier build --no-cache         # Force a clean rebuild without using cache
  phpier build --regenerate       # Regenerate config files before building`,
	RunE: runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Flags
	buildCmd.Flags().BoolVar(&noCache, "no-cache", false, "Do not use cache when building the image")
	buildCmd.Flags().BoolVar(&regenerate, "regenerate", false, "Regenerate configuration files before building")
}

func runBuild(cmd *cobra.Command, args []string) error {
	if !isProjectInitialized() {
		return errors.NewProjectNotInitializedError()
	}

	// Load configurations
	projectCfg, err := config.LoadProjectConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load project config", err)
	}
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load global config", err)
	}

	// Regenerate files only if requested
	if regenerate {
		logrus.Infof("🔄 Regenerating configuration files...")
		engine := templates.NewEngine()
		if err := generator.GenerateProjectFiles(engine, projectCfg, globalCfg); err != nil {
			return errors.WrapError(errors.ErrorTypeUnknown, "Failed to generate project files", err)
		}
	}

	// Create Docker Compose manager
	composeManager, err := docker.NewProjectComposeManager(projectCfg, globalCfg)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client", err)
	}

	// Build the app container
	logrus.Infof("🔨 Building app container...")
	if noCache {
		logrus.Infof("♻️  Using --no-cache flag for clean rebuild")
	}

	if err := composeManager.Build(noCache, "app"); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to build app container", err)
	}

	logrus.Infof("✅ App container built successfully!")
	return nil
}
