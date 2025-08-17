package cmd

import (
	"context"
	"fmt"
	"strings"

	"phpier/internal/config"
	"phpier/internal/docker"

	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage database services",
	Long: `Manage database services (MySQL, PostgreSQL, MariaDB).
	
This command provides database service management capabilities including:
- Enable/disable individual database services
- View status of all database services
- View credentials for enabled services
- List all available database services

Multiple database services can run simultaneously.`,
}

// dbListCmd represents the db list command
var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all database services and their enabled/disabled status",
	Long: `List all available database services and show whether they are enabled or disabled.

This shows the configuration status, not the running status.`,
	RunE: runDbList,
}

// dbStatusCmd represents the db status command
var dbStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show running status of enabled database services with ports",
	Long: `Show the running status of all enabled database services along with their ports.

This command checks if the Docker containers are actually running and displays
their connection information.`,
	RunE: runDbStatus,
}

// dbCredentialsCmd represents the db credentials command
var dbCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Show database credentials for enabled services",
	Long: `Show database connection credentials for all enabled database services.

This displays the host, port, username, and password information needed
to connect to each enabled database service.`,
	RunE: runDbCredentials,
}

// dbEnableCmd represents the db enable command
var dbEnableCmd = &cobra.Command{
	Use:   "enable [mysql|postgresql|mariadb]",
	Short: "Enable a database service",
	Long: `Enable a specific database service (MySQL, PostgreSQL, or MariaDB).

The service will be included in the global docker-compose configuration
and started with the next 'phpier global up' command.`,
	Args: cobra.ExactArgs(1),
	RunE: runDbEnable,
}

// dbDisableCmd represents the db disable command
var dbDisableCmd = &cobra.Command{
	Use:   "disable [mysql|postgresql|mariadb]",
	Short: "Disable a database service",
	Long: `Disable a specific database service (MySQL, PostgreSQL, or MariaDB).

The service will be removed from the global docker-compose configuration
and stopped if currently running.`,
	Args: cobra.ExactArgs(1),
	RunE: runDbDisable,
}

func runDbList(cmd *cobra.Command, args []string) error {
	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	fmt.Println("Database Services:")
	fmt.Printf("✓ MySQL      [%s]\n", getEnabledStatus(globalConfig.Services.Databases.MySQL.Enabled))
	fmt.Printf("✓ PostgreSQL [%s]\n", getEnabledStatus(globalConfig.Services.Databases.PostgreSQL.Enabled))
	fmt.Printf("✓ MariaDB    [%s]\n", getEnabledStatus(globalConfig.Services.Databases.MariaDB.Enabled))

	return nil
}

func runDbStatus(cmd *cobra.Command, args []string) error {
	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	fmt.Println("Database Services Status:")

	// Check MySQL
	if globalConfig.Services.Databases.MySQL.Enabled {
		running := checkDatabaseRunning("mysql")
		fmt.Printf("✓ MySQL      [enabled]  [%s]   localhost:%d\n",
			getRunningStatus(running), globalConfig.Services.Databases.MySQL.Port)
	} else {
		fmt.Printf("✗ MySQL      [disabled] [stopped]   localhost:%d\n",
			globalConfig.Services.Databases.MySQL.Port)
	}

	// Check PostgreSQL
	if globalConfig.Services.Databases.PostgreSQL.Enabled {
		running := checkDatabaseRunning("postgres")
		fmt.Printf("✓ PostgreSQL [enabled]  [%s]   localhost:%d\n",
			getRunningStatus(running), globalConfig.Services.Databases.PostgreSQL.Port)
	} else {
		fmt.Printf("✗ PostgreSQL [disabled] [stopped]   localhost:%d\n",
			globalConfig.Services.Databases.PostgreSQL.Port)
	}

	// Check MariaDB
	if globalConfig.Services.Databases.MariaDB.Enabled {
		running := checkDatabaseRunning("mariadb")
		fmt.Printf("✓ MariaDB    [enabled]  [%s]   localhost:%d\n",
			getRunningStatus(running), globalConfig.Services.Databases.MariaDB.Port)
	} else {
		fmt.Printf("✗ MariaDB    [disabled] [stopped]   localhost:%d\n",
			globalConfig.Services.Databases.MariaDB.Port)
	}

	return nil
}

func runDbCredentials(cmd *cobra.Command, args []string) error {
	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	enabled := globalConfig.GetEnabledDatabases()
	if len(enabled) == 0 {
		fmt.Println("No database services are currently enabled.")
		fmt.Println("Use 'phpier global db enable <service>' to enable a database service.")
		return nil
	}

	fmt.Println("Database Credentials (enabled services only):")
	fmt.Println()

	for dbType, dbConfig := range enabled {
		fmt.Printf("%s:\n", strings.Title(dbType))
		fmt.Printf("  Host:     localhost:%d\n", dbConfig.Port)
		fmt.Printf("  Username: %s\n", dbConfig.Username)
		fmt.Printf("  Password: %s\n", dbConfig.Password)
		fmt.Println()
	}

	return nil
}

func runDbEnable(cmd *cobra.Command, args []string) error {
	dbType := strings.ToLower(args[0])

	if !isValidDatabaseType(dbType) {
		return fmt.Errorf("invalid database type '%s'. Valid options: mysql, postgresql, mariadb", dbType)
	}

	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	// Check if already enabled
	if globalConfig.IsDatabaseEnabled(dbType) {
		fmt.Printf("%s is already enabled.\n", strings.Title(dbType))
		return nil
	}

	// Enable the service with updated defaults
	switch dbType {
	case "mysql":
		globalConfig.Services.Databases.MySQL.Enabled = true
		// Apply updated defaults if not already set
		if globalConfig.Services.Databases.MySQL.Username == "root" {
			globalConfig.Services.Databases.MySQL.Username = "phpier"
			globalConfig.Services.Databases.MySQL.Password = "phpier"
			globalConfig.Services.Databases.MySQL.Database = "phpier"
		}
	case "postgresql":
		globalConfig.Services.Databases.PostgreSQL.Enabled = true
	case "mariadb":
		globalConfig.Services.Databases.MariaDB.Enabled = true
	}

	// Save the updated configuration
	if err := config.SaveGlobalConfig(globalConfig); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}

	fmt.Printf("✓ %s has been enabled.\n", strings.Title(dbType))
	fmt.Println("Run 'phpier global up' to start the service.")

	return nil
}

func runDbDisable(cmd *cobra.Command, args []string) error {
	dbType := strings.ToLower(args[0])

	if !isValidDatabaseType(dbType) {
		return fmt.Errorf("invalid database type '%s'. Valid options: mysql, postgresql, mariadb", dbType)
	}

	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	// Check if already disabled
	if !globalConfig.IsDatabaseEnabled(dbType) {
		fmt.Printf("%s is already disabled.\n", strings.Title(dbType))
		return nil
	}

	// Disable the service
	switch dbType {
	case "mysql":
		globalConfig.Services.Databases.MySQL.Enabled = false
	case "postgresql":
		globalConfig.Services.Databases.PostgreSQL.Enabled = false
	case "mariadb":
		globalConfig.Services.Databases.MariaDB.Enabled = false
	}

	// Save the updated configuration
	if err := config.SaveGlobalConfig(globalConfig); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}

	fmt.Printf("✗ %s has been disabled.\n", strings.Title(dbType))
	fmt.Println("Run 'phpier global up' to apply changes (the service will be stopped).")

	return nil
}

// Helper functions

func getEnabledStatus(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}

func getRunningStatus(running bool) string {
	if running {
		return "running"
	}
	return "stopped"
}

func checkDatabaseRunning(serviceName string) bool {
	dockerClient, err := docker.NewClient()
	if err != nil {
		return false
	}
	defer dockerClient.Close()

	ctx := context.Background()
	containerID, err := dockerClient.GetContainerID("phpier", serviceName)
	if err != nil {
		return false
	}

	running, err := dockerClient.IsContainerRunningByID(ctx, containerID)
	if err != nil {
		return false
	}

	return running
}

func isValidDatabaseType(dbType string) bool {
	validTypes := []string{"mysql", "postgresql", "mariadb"}
	for _, valid := range validTypes {
		if dbType == valid {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbListCmd)
	dbCmd.AddCommand(dbStatusCmd)
	dbCmd.AddCommand(dbCredentialsCmd)
	dbCmd.AddCommand(dbEnableCmd)
	dbCmd.AddCommand(dbDisableCmd)
}
