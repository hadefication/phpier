package cmd

import (
	"os"
	"os/user"

	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool

	// Version information
	buildVersion = "v1.0.1"
	buildCommit  = "unknown"
	buildDate    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "phpier",
	Short: "A CLI tool to manage PHP development using Docker",
	Long: `PHPier is a CLI tool for managing PHP development environments using Docker.

It supports multiple PHP versions (5.6, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, 8.4) 
with Traefik for folder-based domain routing (<directory>.localhost).

Features:
- Multiple PHP version support
- Docker-based environment isolation  
- Traefik reverse proxy with automatic SSL
- Database options (MySQL, PostgreSQL, MariaDB)
- Caching services (Redis, Memcached)
- Development tools (PHPMyAdmin, Mailpit)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		// Get colored output preference
		colored := !viper.GetBool("no-color")

		// Create error handler
		handler := errors.NewErrorHandler(viper.GetBool("verbose"), colored)

		// Handle the error and exit
		handler.Handle(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.phpier.yml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".phpier" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".phpier")
	}

	// Environment variables
	viper.SetEnvPrefix("PHPIER")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
		}
	}

	// Set WWWUSER environment variable if not already set
	setWWWUserEnvVar()

	// Configure logging
	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    false,
	})
}

// setWWWUserEnvVar sets the WWWUSER environment variable to the current user's UID if not already set
func setWWWUserEnvVar() {
	// Check if WWWUSER is already set
	if wwwuser := os.Getenv("WWWUSER"); wwwuser != "" {
		logrus.Debugf("WWWUSER already set to: %s", wwwuser)
		return
	}

	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		logrus.Debugf("Failed to get current user for WWWUSER: %v", err)
		return
	}

	// Convert UID to string and set as environment variable
	uid := currentUser.Uid
	if err := os.Setenv("WWWUSER", uid); err != nil {
		logrus.Debugf("Failed to set WWWUSER environment variable: %v", err)
		return
	}

	logrus.Debugf("Set WWWUSER environment variable to: %s", uid)
}

// SetVersionInfo sets the version information for the CLI
func SetVersionInfo(version, commit, date string) {
	buildVersion = version
	buildCommit = commit
	buildDate = date
}

// isPhpierProject checks if the current directory contains a phpier project
func isPhpierProject() bool {
	if _, err := os.Stat(".phpier.yml"); os.IsNotExist(err) {
		return false
	}
	// TODO: Also check if .phpier.yml contains phpier.managed=true label
	return true
}
