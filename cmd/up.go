package cmd

import (
	"os"
	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"
	"phpier/internal/generator"
	"phpier/internal/templates"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	detached bool
	build    bool
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the project's app container",
	Long: `Start the project's app container and connect it to the global services.

This command will:
- Start the PHP/Nginx container for the current project.
- Ensure the global services network is available.

Examples:
  phpier up                 # Start app container in the foreground
  phpier up -d              # Start app container in the background (detached)
  phpier up --build         # Rebuild the app container image before starting`,
	RunE: runUp,
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Flags
	upCmd.Flags().BoolVarP(&detached, "detach", "d", false, "Run services in the background")
	upCmd.Flags().BoolVar(&build, "build", false, "Build images before starting services")
}

func runUp(cmd *cobra.Command, args []string) error {
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

	// Regenerate files
	engine := templates.NewEngine()
	if err := generator.GenerateProjectFiles(engine, projectCfg, globalCfg); err != nil {
		return errors.WrapError(errors.ErrorTypeUnknown, "Failed to generate project files", err)
	}

	// Create Docker Compose manager
	composeManager, err := docker.NewProjectComposeManager(projectCfg, globalCfg)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client", err)
	}

	// Build image if requested
	if build {
		logrus.Infof("🔨 Building project image...")
		if err := composeManager.Build(false); err != nil {
			return errors.WrapError(errors.ErrorTypeDockerError, "Failed to build project image", err)
		}
	}

	// Start services
	logrus.Infof("🚀 Starting project container...")
	if err := composeManager.Up(detached); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to start project container", err)
	}

	return nil
}

// isProjectInitialized checks if the current directory has been initialized as a phpier project
func isProjectInitialized() bool {
	if _, err := os.Stat(".phpier.yml"); os.IsNotExist(err) {
		return false
	}
	return true
}
