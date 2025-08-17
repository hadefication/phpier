package cmd

import (
	"phpier/internal/config"
	"phpier/internal/errors"
	"phpier/internal/generator"
	"phpier/internal/templates"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	phpVersion  string
	projectName string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [php-version]",
	Short: "Initialize a new phpier project environment",
	Long: `Initialize a new phpier project environment with the specified PHP version.

This command will:
- Create a .phpier.yml file for project-specific settings (PHP version).
- Generate a Dockerfile for the PHP container.
- Generate a docker-compose.yml to run the app container and connect it to the global services network.

Example:
  phpier init 8.3
  phpier init 7.4 --project-name=my-legacy-app`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags
	initCmd.Flags().StringVarP(&phpVersion, "php-version", "p", "8.3", "PHP version to use")
	initCmd.Flags().StringVar(&projectName, "project-name", "", "Project name (defaults to current directory name)")

	// Bind flags to viper
	viper.BindPFlag("php.version", initCmd.Flags().Lookup("php-version"))
	viper.BindPFlag("docker.project_name", initCmd.Flags().Lookup("project-name"))
}

func runInit(cmd *cobra.Command, args []string) error {
	// Parse PHP version from argument if provided
	if len(args) > 0 {
		phpVersion = args[0]
	}

	// Validate PHP version
	if !config.IsValidPHPVersion(phpVersion) {
		return errors.NewInvalidPHPVersionError(phpVersion, config.PHPVersions)
	}

	// Set project name if not provided
	if projectName == "" {
		projectName = config.GetCurrentDir()
	}

	logrus.Infof("Initializing phpier project '%s'...", projectName)
	logrus.Infof("PHP Version: %s", phpVersion)

	// Ensure global config exists, creating it if it doesn't
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigCorrupted, "Failed to load or create global configuration", err)
	}
	logrus.Infof("Using global network: %s", globalCfg.Network)

	// Create project configuration
	projectCfg := createProjectConfig()

	// Create template engine
	engine := templates.NewEngine()

	// Create directory structure for the project
	if err := generator.CreateProjectDirectories(); err != nil {
		return errors.WrapError(errors.ErrorTypeFileSystemError, "Failed to create project directory structure", err)
	}

	// Generate project files
	if err := generator.GenerateProjectFiles(engine, projectCfg, globalCfg); err != nil {
		return errors.WrapError(errors.ErrorTypeTemplateError, "Failed to generate project configuration files", err)
	}

	// Save project configuration
	if err := config.SaveProjectConfig(projectCfg); err != nil {
		return errors.WrapError(errors.ErrorTypeConfigCorrupted, "Failed to save project configuration file", err)
	}

	logrus.Infof("✅ phpier project initialized successfully!")
	logrus.Infof("📂 Configuration saved to .phpier.yml")
	logrus.Infof("🐳 Docker files generated")
	logrus.Infof("🚀 Run 'phpier global up' to start shared services (if not running)")
	logrus.Infof("🚀 Run 'phpier up' to start your project container")

	return nil
}

func createProjectConfig() *config.ProjectConfig {
	cfg := &config.ProjectConfig{}

	// Top-level configuration
	cfg.Name = projectName
	cfg.PHP = phpVersion

	// Set Node.js version based on PHP version
	if phpVersion == "5.6" {
		cfg.Node = "none" // Skip Node.js for PHP 5.6 due to Debian Stretch compatibility issues
	} else {
		cfg.Node = "lts" // Default to latest LTS version for other PHP versions
	}

	// App configuration
	cfg.App.Volumes = []string{"../:/var/www/html"}
	cfg.App.Environment = []string{"APP_ENV=local", "APP_DEBUG=true"}

	return cfg
}

// getDefaultExtensions returns appropriate extensions for each PHP version
// This is used by the template generator for Dockerfile creation
func getDefaultExtensions(phpVersion string) []string {
	coreExtensions := []string{
		"bcmath", "calendar", "curl", "dom", "exif", "ftp", "gd",
		"intl", "mbstring", "mysqli", "opcache", "pdo", "pdo_mysql",
		"pdo_pgsql", "pgsql", "soap", "sockets", "tokenizer", "xml", "zip",
	}

	// Version-specific adjustments can be made here if needed
	return coreExtensions
}

// getDefaultPHPSettings returns default PHP.ini settings
// This is used by the template generator for php.ini creation
func getDefaultPHPSettings() map[string]string {
	return map[string]string{
		"memory_limit":        "256M",
		"upload_max_filesize": "64M",
		"post_max_size":       "64M",
		"max_execution_time":  "300",
		"max_input_vars":      "3000",
		"display_errors":      "On",
		"error_reporting":     "E_ALL",
		"log_errors":          "On",
		"date.timezone":       "UTC",
	}
}
