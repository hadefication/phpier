package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// ProjectConfig represents the project-specific configuration
type ProjectConfig struct {
	Name string    `mapstructure:"name"`
	PHP  string    `mapstructure:"php"`
	Node string    `mapstructure:"node"`
	App  AppConfig `mapstructure:"app"`
}

// GlobalConfig represents the global configuration (~/.phpier/config.yaml)
type GlobalConfig struct {
	Services ServicesConfig `mapstructure:"services"`
	Traefik  TraefikConfig  `mapstructure:"traefik"`
	Network  string         `mapstructure:"network"`
}

// DockerConfig contains Docker-related configuration for the project
type DockerConfig struct {
	ProjectName string `mapstructure:"project_name"`
}

// PHPConfig contains PHP version configuration for the project
type PHPConfig struct {
	Version string `mapstructure:"version"`
}

// AppConfig contains application container configuration
type AppConfig struct {
	Volumes     []string `mapstructure:"volumes"`     // Volume mappings (default: ["./:/var/www/html"])
	Environment []string `mapstructure:"environment"` // Environment variables (optional)
}

// ServicesConfig contains global service configurations
type ServicesConfig struct {
	Database  DatabaseConfig  `mapstructure:"database"`  // Legacy single database config
	Databases DatabasesConfig `mapstructure:"databases"` // New multi-database config
	Cache     CacheConfig     `mapstructure:"cache"`
	Tools     ToolsConfig     `mapstructure:"tools"`
}

// DatabaseConfig contains global database configuration (legacy)
type DatabaseConfig struct {
	Type    string `mapstructure:"type"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

// DatabasesConfig contains multiple database service configurations
type DatabasesConfig struct {
	MySQL      DatabaseServiceConfig `mapstructure:"mysql"`
	PostgreSQL DatabaseServiceConfig `mapstructure:"postgresql"`
	MariaDB    DatabaseServiceConfig `mapstructure:"mariadb"`
}

// DatabaseServiceConfig contains individual database service configuration
type DatabaseServiceConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Version  string `mapstructure:"version"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// CacheConfig contains global cache service configuration
type CacheConfig struct {
	Redis     CacheServiceConfig `mapstructure:"redis"`
	Memcached CacheServiceConfig `mapstructure:"memcached"`
}

// CacheServiceConfig contains individual cache service config
type CacheServiceConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Port    int  `mapstructure:"port"`
}

// MailpitConfig contains Mailpit configuration
type MailpitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Port    int  `mapstructure:"port"`
}

// ToolsConfig contains global development tools configuration
type ToolsConfig struct {
	PHPMyAdmin bool          `mapstructure:"phpmyadmin"`
	Mailpit    MailpitConfig `mapstructure:"mailpit"`
	PgAdmin    bool          `mapstructure:"pgadmin"`
}

// TraefikConfig contains global Traefik configuration
type TraefikConfig struct {
	Domain  string `mapstructure:"domain"`
	Port    int    `mapstructure:"port"`
	SSLPort int    `mapstructure:"ssl_port"`
}

// PHPVersions contains supported PHP versions
var PHPVersions = []string{"5.6", "7.2", "7.3", "7.4", "8.0", "8.1", "8.2", "8.3", "8.4"}

// DatabaseTypes contains supported database types
var DatabaseTypes = []string{"mysql", "postgresql", "mariadb"}

// LoadProjectConfig loads the project-specific configuration from .phpier.yml labels
func LoadProjectConfig() (*ProjectConfig, error) {
	// Read .phpier.yml to extract project information from labels
	return LoadProjectConfigFromDockerCompose(".phpier.yml")
}

// LoadProjectConfigFromDockerCompose loads project config from .phpier.yml file
func LoadProjectConfigFromDockerCompose(path string) (*ProjectConfig, error) {
	// For now, return a basic config based on current directory and defaults
	// This will be enhanced to read from .phpier.yml labels
	projectName := GetCurrentDir()

	return &ProjectConfig{
		Name: projectName,
		PHP:  "8.3", // Default, will be read from docker-compose.yml labels
		Node: "lts", // Default, will be read from docker-compose.yml labels
		App: AppConfig{
			Volumes:     []string{"./:/var/www/html"},
			Environment: []string{"APP_ENV=local", "APP_DEBUG=true"},
		},
	}, nil
}

// LoadGlobalConfig loads the global configuration from ~/.phpier/config.yaml
func LoadGlobalConfig() (*GlobalConfig, error) {
	globalViper := viper.New()
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	configPath := filepath.Join(home, ".phpier")
	globalViper.SetConfigName("config")
	globalViper.SetConfigType("yaml")
	globalViper.AddConfigPath(configPath)

	setGlobalDefaults(globalViper)

	if err := globalViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create the file if it doesn't exist
			if err := os.MkdirAll(configPath, 0755); err != nil {
				return nil, fmt.Errorf("failed to create global config directory: %w", err)
			}
			if err := globalViper.SafeWriteConfig(); err != nil {
				return nil, fmt.Errorf("failed to write global config file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read global config: %w", err)
		}
	}

	var config GlobalConfig
	if err := globalViper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode global config: %w", err)
	}

	// Migrate legacy configuration if needed
	if migrated := migrateFromLegacyConfig(&config); migrated {
		// Save the migrated configuration
		if err := SaveGlobalConfig(&config); err != nil {
			// Log warning but don't fail - user can still use legacy config
			fmt.Printf("Warning: Failed to save migrated configuration: %v\n", err)
		}
	}

	return &config, nil
}

// CreateProjectConfig creates a project configuration from CLI arguments
func CreateProjectConfig(name, phpVersion, nodeVersion string) *ProjectConfig {
	// Set defaults if not provided
	if name == "" {
		name = GetCurrentDir()
	}
	if phpVersion == "" {
		phpVersion = "8.3"
	}
	if nodeVersion == "" {
		if phpVersion == "5.6" {
			nodeVersion = "none" // Skip Node.js for PHP 5.6
		} else {
			nodeVersion = "lts"
		}
	}

	return &ProjectConfig{
		Name: name,
		PHP:  phpVersion,
		Node: nodeVersion,
		App: AppConfig{
			Volumes:     []string{"./:/var/www/html"},
			Environment: []string{"APP_ENV=local", "APP_DEBUG=true"},
		},
	}
}

// SaveGlobalConfig saves the global configuration to ~/.phpier/config.yaml
func SaveGlobalConfig(config *GlobalConfig) error {
	globalViper := viper.New()
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	configPath := filepath.Join(home, ".phpier")
	globalViper.SetConfigName("config")
	globalViper.SetConfigType("yaml")
	globalViper.AddConfigPath(configPath)

	globalViper.Set("services", config.Services)
	globalViper.Set("traefik", config.Traefik)
	globalViper.Set("network", config.Network)

	if err := globalViper.SafeWriteConfig(); err != nil {
		// If file exists, use WriteConfig to overwrite
		if err := globalViper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to write global config file: %w", err)
		}
	}
	return nil
}

func setGlobalDefaults(v *viper.Viper) {
	v.SetDefault("network", "phpier_global")
	v.SetDefault("traefik.domain", "localhost")
	v.SetDefault("traefik.port", 80)
	v.SetDefault("traefik.ssl_port", 443)

	// Legacy single database config (for backwards compatibility)
	v.SetDefault("services.database.type", "mysql")
	v.SetDefault("services.database.version", "8.0")
	v.SetDefault("services.database.port", 3306)

	// New multi-database config
	v.SetDefault("services.databases.mysql.enabled", true)
	v.SetDefault("services.databases.mysql.version", "8.0")
	v.SetDefault("services.databases.mysql.port", 3306)
	v.SetDefault("services.databases.mysql.username", "phpier")
	v.SetDefault("services.databases.mysql.password", "phpier")
	v.SetDefault("services.databases.mysql.database", "phpier")

	v.SetDefault("services.databases.postgresql.enabled", false)
	v.SetDefault("services.databases.postgresql.version", "15")
	v.SetDefault("services.databases.postgresql.port", 5432)
	v.SetDefault("services.databases.postgresql.username", "phpier")
	v.SetDefault("services.databases.postgresql.password", "phpier")
	v.SetDefault("services.databases.postgresql.database", "phpier")

	v.SetDefault("services.databases.mariadb.enabled", false)
	v.SetDefault("services.databases.mariadb.version", "10.11")
	v.SetDefault("services.databases.mariadb.port", 3307)
	v.SetDefault("services.databases.mariadb.username", "phpier")
	v.SetDefault("services.databases.mariadb.password", "phpier")
	v.SetDefault("services.databases.mariadb.database", "phpier")

	v.SetDefault("services.cache.redis.enabled", true)
	v.SetDefault("services.cache.redis.port", 6379)
	v.SetDefault("services.cache.memcached.enabled", false)
	v.SetDefault("services.cache.memcached.port", 11211)

	v.SetDefault("services.tools.phpmyadmin", true)
	v.SetDefault("services.tools.mailpit.enabled", true)
	v.SetDefault("services.tools.mailpit.port", 1025)
	v.SetDefault("services.tools.pgadmin", false)
}

// GetCurrentDir returns the current directory name for domain generation
func GetCurrentDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		return "phpier"
	}
	return filepath.Base(pwd)
}

// IsValidPHPVersion checks if a PHP version is supported.
func IsValidPHPVersion(version string) bool {
	for _, v := range PHPVersions {
		if v == version {
			return true
		}
	}
	return false
}

// GetEnabledDatabases returns a list of enabled database services
func (c *GlobalConfig) GetEnabledDatabases() map[string]DatabaseServiceConfig {
	enabled := make(map[string]DatabaseServiceConfig)

	if c.Services.Databases.MySQL.Enabled {
		enabled["mysql"] = c.Services.Databases.MySQL
	}
	if c.Services.Databases.PostgreSQL.Enabled {
		enabled["postgresql"] = c.Services.Databases.PostgreSQL
	}
	if c.Services.Databases.MariaDB.Enabled {
		enabled["mariadb"] = c.Services.Databases.MariaDB
	}

	return enabled
}

// GetDatabaseService returns the configuration for a specific database service
func (c *GlobalConfig) GetDatabaseService(dbType string) (DatabaseServiceConfig, bool) {
	switch dbType {
	case "mysql":
		return c.Services.Databases.MySQL, true
	case "postgresql":
		return c.Services.Databases.PostgreSQL, true
	case "mariadb":
		return c.Services.Databases.MariaDB, true
	default:
		return DatabaseServiceConfig{}, false
	}
}

// IsDatabaseEnabled checks if a specific database service is enabled
func (c *GlobalConfig) IsDatabaseEnabled(dbType string) bool {
	config, exists := c.GetDatabaseService(dbType)
	return exists && config.Enabled
}

// migrateFromLegacyConfig migrates legacy single database configuration to new multi-database format
func migrateFromLegacyConfig(config *GlobalConfig) bool {
	// Check if we have legacy config and no new config
	if config.Services.Database.Type != "" && !hasNewDatabaseConfig(config) {
		fmt.Println("Migrating legacy database configuration to new multi-database format...")

		// Enable the legacy database type in the new format
		switch config.Services.Database.Type {
		case "mysql":
			config.Services.Databases.MySQL.Enabled = true
			config.Services.Databases.MySQL.Version = config.Services.Database.Version
			config.Services.Databases.MySQL.Port = config.Services.Database.Port
		case "postgresql":
			config.Services.Databases.PostgreSQL.Enabled = true
			config.Services.Databases.PostgreSQL.Version = config.Services.Database.Version
			config.Services.Databases.PostgreSQL.Port = config.Services.Database.Port
		case "mariadb":
			config.Services.Databases.MariaDB.Enabled = true
			config.Services.Databases.MariaDB.Version = config.Services.Database.Version
			config.Services.Databases.MariaDB.Port = config.Services.Database.Port
		}

		fmt.Printf("✓ Migrated %s database configuration\n", config.Services.Database.Type)
		fmt.Println("Use 'phpier global db list' to see enabled database services")
		return true
	}
	return false
}

// hasNewDatabaseConfig checks if the new multi-database configuration is present
func hasNewDatabaseConfig(config *GlobalConfig) bool {
	return config.Services.Databases.MySQL.Enabled ||
		config.Services.Databases.PostgreSQL.Enabled ||
		config.Services.Databases.MariaDB.Enabled
}
