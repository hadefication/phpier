package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"phpier/internal/errors"
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

// ProjectInfo contains information about a discovered phpier project
type ProjectInfo struct {
	Name string
	Path string
}

// FindProjectByName searches for a phpier project by name using Docker and filesystem discovery
func FindProjectByName(projectName string) (*ProjectInfo, error) {
	// Get projects from both sources
	dockerProjects, dockerErr := DiscoverProjectsFromDocker()
	filesystemProjects, filesystemErr := DiscoverProjectsFromFilesystem()

	// Combine projects, preferring filesystem entries that have valid paths
	allProjects := make(map[string]ProjectInfo)
	
	// Add Docker projects first
	if dockerErr == nil {
		for _, project := range dockerProjects {
			allProjects[project.Name] = project
		}
	}
	
	// Add filesystem projects, replacing Docker entries that lack paths
	if filesystemErr == nil {
		for _, project := range filesystemProjects {
			if existing, exists := allProjects[project.Name]; !exists || existing.Path == "" {
				allProjects[project.Name] = project
			}
		}
	}

	// Convert back to slice for findProjectInList
	var combinedProjects []ProjectInfo
	for _, project := range allProjects {
		combinedProjects = append(combinedProjects, project)
	}

	if len(combinedProjects) == 0 {
		if dockerErr != nil && filesystemErr != nil {
			return nil, fmt.Errorf("failed to discover projects: docker error: %v, filesystem error: %v", dockerErr, filesystemErr)
		}
		return nil, errors.NewProjectNotFoundError(projectName)
	}

	return findProjectInList(projectName, combinedProjects)
}

// findProjectInList searches for a project by name in a given list of projects
func findProjectInList(projectName string, projects []ProjectInfo) (*ProjectInfo, error) {
	var matches []ProjectInfo
	for _, project := range projects {
		if project.Name == projectName {
			matches = append(matches, project)
		}
	}

	if len(matches) == 0 {
		return nil, errors.NewProjectNotFoundError(projectName)
	}

	if len(matches) > 1 {
		var paths []string
		for _, match := range matches {
			paths = append(paths, match.Path)
		}
		return nil, errors.NewMultipleProjectsFoundError(projectName, paths)
	}

	return &matches[0], nil
}

// DiscoverProjectsFromDocker discovers phpier projects by scanning Docker images with phpier- prefix
func DiscoverProjectsFromDocker() ([]ProjectInfo, error) {
	// Check if Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		return []ProjectInfo{}, nil // Docker not available, return empty list
	}

	// Get all phpier- prefixed images
	cmd := exec.Command("docker", "images", "--filter", "reference=phpier-*", "--format", "{{.Repository}}:{{.Tag}}")
	output, err := cmd.Output()
	if err != nil {
		return []ProjectInfo{}, nil // Docker command failed, return empty list
	}

	if strings.TrimSpace(string(output)) == "" {
		return []ProjectInfo{}, nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	projectMap := make(map[string]ProjectInfo)

	for _, line := range lines {
		imageName := strings.TrimSpace(line)
		if imageName == "" || !strings.HasPrefix(imageName, "phpier-") {
			continue
		}

		// Extract project name from image name (remove phpier- prefix and :tag)
		projectName := strings.TrimPrefix(imageName, "phpier-")
		if colonIndex := strings.Index(projectName, ":"); colonIndex != -1 {
			projectName = projectName[:colonIndex]
		}

		if projectName == "" {
			continue
		}

		// Check if this project has containers and get working directory
		workingDir := ""
		
		// Try to get container info for this project
		containerCmd := exec.Command("docker", "ps", "-a",
			"--filter", fmt.Sprintf("ancestor=%s", imageName),
			"--format", "{{.Status}}\t{{.Label \"com.docker.compose.working-dir\"}}")
		
		if containerOutput, err := containerCmd.Output(); err == nil {
			containerInfo := strings.TrimSpace(string(containerOutput))
			if containerInfo != "" {
				parts := strings.Split(containerInfo, "\t")
				if len(parts) > 1 {
					workingDir = strings.TrimSpace(parts[1])
				}
			}
		}

		// Try to find working directory from Docker Compose project label if not found
		if workingDir == "" {
			projectCmd := exec.Command("docker", "ps", "-a",
				"--filter", fmt.Sprintf("label=com.docker.compose.project=%s", projectName),
				"--format", "{{.Label \"com.docker.compose.working-dir\"}}")
			
			if projectOutput, err := projectCmd.Output(); err == nil {
				workingDir = strings.TrimSpace(string(projectOutput))
			}
		}

		projectMap[projectName] = ProjectInfo{
			Name: projectName,
			Path: workingDir,
		}
	}

	// Convert map to slice
	var result []ProjectInfo
	for _, project := range projectMap {
		result = append(result, project)
	}

	return result, nil
}

// DiscoverProjectsFromFilesystem scans the filesystem for phpier projects
func DiscoverProjectsFromFilesystem() ([]ProjectInfo, error) {
	return DiscoverProjects()
}

// DiscoverProjects scans the filesystem for phpier projects (legacy function)
func DiscoverProjects() ([]ProjectInfo, error) {
	var projects []ProjectInfo

	// Common search paths
	searchPaths := []string{
		".", // Current directory
	}

	// Add user home directory subdirectories if accessible
	if home, err := os.UserHomeDir(); err == nil {
		commonDirs := []string{
			filepath.Join(home, "projects"),
			filepath.Join(home, "code"),
			filepath.Join(home, "dev"),
			filepath.Join(home, "Sites"),
			filepath.Join(home, "workspace"),
			filepath.Join(home, "Development"),
		}
		searchPaths = append(searchPaths, commonDirs...)
	}

	// Add common development directories
	commonPaths := []string{
		"/var/www",
		"/srv",
		"/opt/projects",
	}
	searchPaths = append(searchPaths, commonPaths...)

	for _, searchPath := range searchPaths {
		if _, err := os.Stat(searchPath); os.IsNotExist(err) {
			continue
		}

		// Search for .phpier.yml files recursively (up to 3 levels deep)
		err := scanForProjects(searchPath, 0, 3, &projects)
		if err != nil {
			// Continue searching other paths even if one fails
			continue
		}
	}

	return projects, nil
}

// scanForProjects recursively scans for .phpier.yml files
func scanForProjects(dir string, currentDepth, maxDepth int, projects *[]ProjectInfo) error {
	if currentDepth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Check if current directory has .phpier.yml
	for _, entry := range entries {
		if entry.Name() == ".phpier.yml" {
			projectInfo, err := extractProjectInfo(dir)
			if err != nil {
				continue // Skip invalid projects
			}
			*projects = append(*projects, *projectInfo)
			break // Found project file, don't recurse further in this directory
		}
	}

	// If no .phpier.yml found, recurse into subdirectories
	if currentDepth < maxDepth {
		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				subDir := filepath.Join(dir, entry.Name())
				scanForProjects(subDir, currentDepth+1, maxDepth, projects)
			}
		}
	}

	return nil
}

// extractProjectInfo extracts project information from a directory containing .phpier.yml
func extractProjectInfo(projectPath string) (*ProjectInfo, error) {
	configPath := filepath.Join(projectPath, ".phpier.yml")

	// For now, use directory name as project name
	// In the future, this could read the actual project name from .phpier.yml labels
	projectName := filepath.Base(projectPath)

	// Verify the file is readable
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no .phpier.yml file found in %s", projectPath)
	}

	return &ProjectInfo{
		Name: projectName,
		Path: projectPath,
	}, nil
}

// LoadProjectConfigFromPath loads project configuration from a specific path
func LoadProjectConfigFromPath(projectPath string) (*ProjectConfig, error) {
	// For now, return a basic config based on directory name
	// This will be enhanced to read from .phpier.yml labels when that feature is implemented
	projectName := filepath.Base(projectPath)

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
