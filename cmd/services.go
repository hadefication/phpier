package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"phpier/internal/config"
	"phpier/internal/display"
	"phpier/internal/docker"
	"phpier/internal/errors"
)

var (
	servicesProjectFilter string
	servicesTypeFilter    string
	servicesStatusFilter  string
	servicesJSONOutput    bool
	servicesVerbose       bool
	servicesShowPorts     bool
	servicesShowURLs      bool
)

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Show status of all phpier services",
	Long: `Display comprehensive information about all running phpier services.

This command shows the status of global phpier services (like Traefik) and 
project-specific services (app containers, databases, caching, development tools).

The output includes:
- Service name and status (running/stopped/exited)
- Container uptime and health information
- Port mappings and exposed services
- Service URLs for web-accessible tools
- Resource information in verbose mode

Examples:
  phpier services                           # Show all phpier services
  phpier services --project myapp           # Show services for specific project
  phpier services --type app                # Show only app containers
  phpier services --status running          # Show only running services
  phpier services --json                    # Output in JSON format
  phpier services --verbose                 # Show detailed information
  phpier services --project myapp --verbose # Detailed view of project services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServices(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)

	// Filter flags
	servicesCmd.Flags().StringVarP(&servicesProjectFilter, "project", "p", "", "Filter by project name")
	servicesCmd.Flags().StringVarP(&servicesTypeFilter, "type", "t", "", "Filter by service type (app, db, cache, proxy, tools)")
	servicesCmd.Flags().StringVarP(&servicesStatusFilter, "status", "s", "", "Filter by status (running, stopped, exited)")

	// Output format flags
	servicesCmd.Flags().BoolVar(&servicesJSONOutput, "json", false, "Output in JSON format")
	servicesCmd.Flags().BoolVarP(&servicesVerbose, "verbose", "v", false, "Show detailed service information")

	// Display options flags
	servicesCmd.Flags().BoolVar(&servicesShowPorts, "ports", true, "Show port mappings")
	servicesCmd.Flags().BoolVar(&servicesShowURLs, "urls", true, "Show service URLs")
}

func runServices(cmd *cobra.Command, args []string) error {
	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	ctx := context.Background()

	// Prepare filter
	filter := &docker.ServicesFilter{}

	if servicesProjectFilter != "" {
		filter.Project = servicesProjectFilter
	}

	if servicesStatusFilter != "" {
		filter.Status = servicesStatusFilter
	}

	// Apply service type filter logic
	if servicesTypeFilter != "" {
		// We'll filter after getting all services since type filtering
		// requires analyzing service names and labels
	}

	// Get services
	services, err := dockerClient.GetPhpierServices(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to get phpier services: %w", err)
	}

	// Apply type filter if specified
	if servicesTypeFilter != "" {
		services = filterServicesByType(services, servicesTypeFilter)
	}

	// If no project filter specified and we're in a phpier project directory,
	// try to detect the current project
	if servicesProjectFilter == "" {
		currentProject := detectCurrentProject()
		if currentProject != "" {
			// If we have services and current project is detected,
			// separate them appropriately
			if len(services) > 0 {
				// This will be handled in the display logic
			}
		}
	}

	// Output in JSON format if requested
	if servicesJSONOutput {
		return outputServicesJSON(services)
	}

	// Output in table format
	return outputServicesTable(services)
}

// filterServicesByType filters services by their type
func filterServicesByType(services []docker.ServiceInfo, serviceType string) []docker.ServiceInfo {
	var filtered []docker.ServiceInfo

	for _, service := range services {
		if matchesServiceType(service, serviceType) {
			filtered = append(filtered, service)
		}
	}

	return filtered
}

// matchesServiceType checks if a service matches the specified type
func matchesServiceType(service docker.ServiceInfo, serviceType string) bool {
	serviceName := strings.ToLower(service.Service)
	containerName := strings.ToLower(service.Name)

	switch strings.ToLower(serviceType) {
	case "app":
		return serviceName == "app" || strings.Contains(containerName, "-app-")

	case "db", "database":
		dbServices := []string{"mysql", "postgres", "mariadb", "postgresql"}
		for _, db := range dbServices {
			if serviceName == db || strings.Contains(containerName, "-"+db+"-") {
				return true
			}
		}
		return false

	case "cache":
		cacheServices := []string{"redis", "valkey", "memcached"}
		for _, cache := range cacheServices {
			if serviceName == cache || strings.Contains(containerName, "-"+cache+"-") {
				return true
			}
		}
		return false

	case "proxy":
		return serviceName == "traefik" || strings.Contains(containerName, "traefik")

	case "tools":
		toolServices := []string{"phpmyadmin", "mailpit", "adminer"}
		for _, tool := range toolServices {
			if serviceName == tool || strings.Contains(containerName, "-"+tool+"-") {
				return true
			}
		}
		return false

	default:
		return false
	}
}

// detectCurrentProject tries to detect the current project from .phpier.yml
func detectCurrentProject() string {
	// Look for .phpier.yml in current directory
	if _, err := os.Stat(".phpier.yml"); err == nil {
		cfg, err := config.LoadProjectConfig()
		if err == nil && cfg.Name != "" {
			return cfg.Name
		}
	}

	// Try to detect from current directory name
	currentDir, err := os.Getwd()
	if err == nil {
		baseName := filepath.Base(currentDir)
		// Basic validation - avoid obviously non-project names
		if baseName != "." && baseName != "/" && len(baseName) > 0 {
			return baseName
		}
	}

	return ""
}

// outputServicesJSON outputs services in JSON format
func outputServicesJSON(services []docker.ServiceInfo) error {
	output := map[string]interface{}{
		"services": services,
		"count":    len(services),
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal services to JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

// outputServicesTable outputs services in table format
func outputServicesTable(services []docker.ServiceInfo) error {
	// Prepare display options
	displayOptions := display.TableOptions{
		ShowHeaders: true,
		ColorOutput: !viper.GetBool("no-color"),
		MaxWidth:    120,
	}

	// Prepare table configuration
	tableConfig := display.ServicesTableConfig{
		ShowPorts:    servicesShowPorts,
		ShowMounts:   servicesVerbose,
		ShowNetworks: servicesVerbose,
		ShowURLs:     servicesShowURLs,
		Verbose:      servicesVerbose,
	}

	// Render table or verbose output
	var output string
	if servicesVerbose {
		output = display.RenderVerboseServicesInfo(services, displayOptions)
	} else {
		output = display.RenderServicesTable(services, tableConfig, displayOptions)
	}

	fmt.Println(output)

	// Show helpful information if no services found
	if len(services) == 0 {
		fmt.Println()
		fmt.Println("To start phpier services:")
		fmt.Println("  phpier global up          # Start global services")
		fmt.Println("  phpier init 8.3           # Initialize a new project")
		fmt.Println("  phpier up                 # Start project services")
	}

	return nil
}

// Additional helper functions for enhanced functionality

// getProjectConfigPath returns the path to project configuration
func getProjectConfigPath() string {
	if _, err := os.Stat(".phpier.yml"); err == nil {
		return ".phpier.yml"
	}
	return ""
}

// isInPhpierProject checks if current directory is a phpier project
func isInPhpierProject() bool {
	return getProjectConfigPath() != ""
}

// validateServiceType validates the service type filter
func validateServiceType(serviceType string) error {
	validTypes := []string{"app", "db", "database", "cache", "proxy", "tools"}

	for _, validType := range validTypes {
		if strings.EqualFold(serviceType, validType) {
			return nil
		}
	}

	return errors.NewInvalidArgumentsError(fmt.Sprintf("invalid service type '%s'. Valid types: %s",
		serviceType, strings.Join(validTypes, ", ")))
}

// validateStatusFilter validates the status filter
func validateStatusFilter(status string) error {
	validStatuses := []string{"running", "stopped", "exited", "paused", "restarting", "created"}

	for _, validStatus := range validStatuses {
		if strings.EqualFold(status, validStatus) {
			return nil
		}
	}

	return errors.NewInvalidArgumentsError(fmt.Sprintf("invalid status '%s'. Valid statuses: %s",
		status, strings.Join(validStatuses, ", ")))
}
