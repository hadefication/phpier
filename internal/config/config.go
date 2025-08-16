package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// ProjectConfig represents the project-specific configuration (.phpier.yml)
type ProjectConfig struct {
	Name string    `mapstructure:"name"`
	PHP  string    `mapstructure:"php"`
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
	Database DatabaseConfig `mapstructure:"database"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Tools    ToolsConfig    `mapstructure:"tools"`
}

// DatabaseConfig contains global database configuration
type DatabaseConfig struct {
	Type    string `mapstructure:"type"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
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
var PHPVersions = []string{"5.6", "7.3", "7.4", "8.0", "8.1", "8.2", "8.3", "8.4"}

// DatabaseTypes contains supported database types
var DatabaseTypes = []string{"mysql", "postgresql", "mariadb"}

// LoadProjectConfig loads the project-specific configuration from .phpier.yml
func LoadProjectConfig() (*ProjectConfig, error) {
	projectViper := viper.New()
	projectViper.SetConfigFile(".phpier.yml")

	if err := projectViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read project config: %w", err)
	}

	var config ProjectConfig
	if err := projectViper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode project config: %w", err)
	}

	return &config, nil
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

	return &config, nil
}

// SaveProjectConfig saves the project-specific configuration to .phpier.yml
func SaveProjectConfig(config *ProjectConfig) error {
	projectViper := viper.New()
	projectViper.SetConfigFile(".phpier.yml")

	projectViper.Set("name", config.Name)
	projectViper.Set("php", config.PHP)
	projectViper.Set("app", config.App)

	if err := projectViper.SafeWriteConfig(); err != nil {
		// If file exists, use WriteConfig to overwrite
		if err := projectViper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to write project config file: %w", err)
		}
	}
	return nil
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
		return fmt.Errorf("failed to write global config file: %w", err)
	}
	return nil
}

func setGlobalDefaults(v *viper.Viper) {
	v.SetDefault("network", "phpier_global")
	v.SetDefault("traefik.domain", "localhost")
	v.SetDefault("traefik.port", 80)
	v.SetDefault("traefik.ssl_port", 443)

	v.SetDefault("services.database.type", "mysql")
	v.SetDefault("services.database.version", "8.0")
	v.SetDefault("services.database.port", 3306)

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
