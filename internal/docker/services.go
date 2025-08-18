package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ServiceInfo represents information about a Docker service/container
type ServiceInfo struct {
	Name       string            `json:"name"`
	Image      string            `json:"image"`
	Status     string            `json:"status"`
	State      string            `json:"state"`
	Health     string            `json:"health"`
	Ports      []PortMapping     `json:"ports"`
	Created    time.Time         `json:"created"`
	StartedAt  time.Time         `json:"started_at"`
	Uptime     string            `json:"uptime"`
	Project    string            `json:"project"`
	Service    string            `json:"service"`
	Labels     map[string]string `json:"labels"`
	Mounts     []MountInfo       `json:"mounts"`
	Networks   []string          `json:"networks"`
	ServiceURL string            `json:"service_url,omitempty"`
}

// PortMapping represents a port mapping
type PortMapping struct {
	ContainerPort string `json:"container_port"`
	HostPort      string `json:"host_port"`
	Protocol      string `json:"protocol"`
}

// MountInfo represents mount information
type MountInfo struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Type        string `json:"type"`
	Mode        string `json:"mode"`
}

// ServicesFilter represents filtering options for services
type ServicesFilter struct {
	Project     string
	ServiceType string
	Status      string
}

// GetPhpierServices returns information about all phpier-related services
func (c *Client) GetPhpierServices(ctx context.Context, filter *ServicesFilter) ([]ServiceInfo, error) {
	var services []ServiceInfo

	// Get all containers with phpier labels or naming patterns
	containers, err := c.getContainersByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get containers: %w", err)
	}

	for _, container := range containers {
		serviceInfo, err := c.getServiceInfo(ctx, container)
		if err != nil {
			logrus.Warnf("Failed to get info for container %s: %v", container, err)
			continue
		}
		services = append(services, serviceInfo)
	}

	// Sort services by project, then by service type
	sort.Slice(services, func(i, j int) bool {
		if services[i].Project != services[j].Project {
			return services[i].Project < services[j].Project
		}
		return services[i].Service < services[j].Service
	})

	return services, nil
}

// getContainersByFilter gets containers matching the filter criteria
func (c *Client) getContainersByFilter(ctx context.Context, filter *ServicesFilter) ([]string, error) {
	args := []string{"ps", "-a", "--format", "{{.Names}}"}

	// Add filter for phpier containers
	args = append(args, "--filter", "label=com.docker.compose.project")

	// Add project filter if specified
	if filter != nil && filter.Project != "" {
		args = append(args, "--filter", fmt.Sprintf("label=com.docker.compose.project=%s", filter.Project))
	}

	// Add status filter if specified
	if filter != nil && filter.Status != "" {
		args = append(args, "--filter", fmt.Sprintf("status=%s", filter.Status))
	}

	output, err := c.RunCommandOutput("docker", args...)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(output) == "" {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var containers []string

	for _, line := range lines {
		containerName := strings.TrimSpace(line)
		if containerName != "" {
			// Filter for phpier-related containers
			if c.isPhpierContainer(containerName) {
				containers = append(containers, containerName)
			}
		}
	}

	return containers, nil
}

// isPhpierContainer checks if a container is phpier-related
func (c *Client) isPhpierContainer(containerName string) bool {
	// First, exclude known non-phpier patterns
	excludePatterns := []string{
		`^not-phpier-.*`,   // Explicitly not phpier
		`^wordpress-.*`,    // WordPress containers
		`^laravel-.*`,      // Laravel containers
		`^drupal-.*`,       // Drupal containers
		`^magento-.*`,      // Magento containers
	}

	for _, pattern := range excludePatterns {
		matched, err := regexp.MatchString(pattern, containerName)
		if err == nil && matched {
			return false
		}
	}

	phpierPatterns := []string{
		`^phpier-.*`,            // Global phpier containers
		`.*-app-\d+$`,           // Project app containers
		`.*-mysql-\d+$`,         // Project mysql containers
		`.*-postgres-\d+$`,      // Project postgres containers
		`.*-mariadb-\d+$`,       // Project mariadb containers
		`.*-redis-\d+$`,         // Project redis containers
		`.*-valkey-\d+$`,        // Project valkey containers
		`.*-memcached-\d+$`,     // Project memcached containers
		`.*-phpmyadmin-\d+$`,    // Project phpmyadmin containers
		`.*-mailpit-\d+$`,       // Project mailpit containers
	}

	for _, pattern := range phpierPatterns {
		matched, err := regexp.MatchString(pattern, containerName)
		if err == nil && matched {
			return true
		}
	}

	return false
}

// getServiceInfo gets detailed information about a specific container
func (c *Client) getServiceInfo(ctx context.Context, containerName string) (ServiceInfo, error) {
	// Get container inspect information
	inspectOutput, err := c.RunCommandOutput("docker", "inspect", containerName)
	if err != nil {
		return ServiceInfo{}, fmt.Errorf("failed to inspect container %s: %w", containerName, err)
	}

	var containers []map[string]interface{}
	if err := json.Unmarshal([]byte(inspectOutput), &containers); err != nil {
		return ServiceInfo{}, fmt.Errorf("failed to parse inspect output: %w", err)
	}

	if len(containers) == 0 {
		return ServiceInfo{}, fmt.Errorf("no container data found")
	}

	container := containers[0]

	// Extract basic information
	serviceInfo := ServiceInfo{
		Name:   containerName,
		Labels: make(map[string]string),
	}

	// Extract image
	if config, ok := container["Config"].(map[string]interface{}); ok {
		if image, ok := config["Image"].(string); ok {
			serviceInfo.Image = image
		}

		// Extract labels
		if labels, ok := config["Labels"].(map[string]interface{}); ok {
			for k, v := range labels {
				if str, ok := v.(string); ok {
					serviceInfo.Labels[k] = str
				}
			}
		}
	}

	// Extract state information
	if state, ok := container["State"].(map[string]interface{}); ok {
		if status, ok := state["Status"].(string); ok {
			serviceInfo.Status = status
			serviceInfo.State = status
		}

		if health, ok := state["Health"].(map[string]interface{}); ok {
			if healthStatus, ok := health["Status"].(string); ok {
				serviceInfo.Health = healthStatus
			}
		}

		// Extract timestamps
		if createdStr, ok := state["StartedAt"].(string); ok && createdStr != "" {
			if startedAt, err := time.Parse(time.RFC3339Nano, createdStr); err == nil {
				serviceInfo.StartedAt = startedAt
				if serviceInfo.Status == "running" {
					serviceInfo.Uptime = formatUptime(time.Since(startedAt))
				}
			}
		}
	}

	// Extract creation time
	if createdStr, ok := container["Created"].(string); ok && createdStr != "" {
		if created, err := time.Parse(time.RFC3339Nano, createdStr); err == nil {
			serviceInfo.Created = created
		}
	}

	// Extract port mappings
	if networkSettings, ok := container["NetworkSettings"].(map[string]interface{}); ok {
		if ports, ok := networkSettings["Ports"].(map[string]interface{}); ok {
			serviceInfo.Ports = extractPorts(ports)
		}

		// Extract networks
		if networks, ok := networkSettings["Networks"].(map[string]interface{}); ok {
			for networkName := range networks {
				serviceInfo.Networks = append(serviceInfo.Networks, networkName)
			}
		}
	}

	// Extract mounts
	if mounts, ok := container["Mounts"].([]interface{}); ok {
		serviceInfo.Mounts = extractMounts(mounts)
	}

	// Extract project and service information from labels
	if project, ok := serviceInfo.Labels["com.docker.compose.project"]; ok {
		serviceInfo.Project = project
	}
	if service, ok := serviceInfo.Labels["com.docker.compose.service"]; ok {
		serviceInfo.Service = service
	}

	// Generate service URL if applicable
	serviceInfo.ServiceURL = c.generateServiceURL(serviceInfo)

	return serviceInfo, nil
}

// extractPorts extracts port mappings from Docker inspect output
func extractPorts(ports map[string]interface{}) []PortMapping {
	var portMappings []PortMapping

	for containerPort, bindings := range ports {
		if bindingList, ok := bindings.([]interface{}); ok && len(bindingList) > 0 {
			for _, binding := range bindingList {
				if bindingMap, ok := binding.(map[string]interface{}); ok {
					hostPort, _ := bindingMap["HostPort"].(string)
					hostIP, _ := bindingMap["HostIp"].(string)

					parts := strings.Split(containerPort, "/")
					port := parts[0]
					protocol := "tcp"
					if len(parts) > 1 {
						protocol = parts[1]
					}

					mapping := PortMapping{
						ContainerPort: port,
						HostPort:      hostPort,
						Protocol:      protocol,
					}

					if hostIP != "" && hostIP != "0.0.0.0" {
						mapping.HostPort = fmt.Sprintf("%s:%s", hostIP, hostPort)
					}

					portMappings = append(portMappings, mapping)
				}
			}
		}
	}

	return portMappings
}

// extractMounts extracts mount information from Docker inspect output
func extractMounts(mounts []interface{}) []MountInfo {
	var mountInfo []MountInfo

	for _, mount := range mounts {
		if mountMap, ok := mount.(map[string]interface{}); ok {
			info := MountInfo{}

			if source, ok := mountMap["Source"].(string); ok {
				info.Source = source
			}
			if destination, ok := mountMap["Destination"].(string); ok {
				info.Destination = destination
			}
			if mountType, ok := mountMap["Type"].(string); ok {
				info.Type = mountType
			}
			if mode, ok := mountMap["Mode"].(string); ok {
				info.Mode = mode
			}

			mountInfo = append(mountInfo, info)
		}
	}

	return mountInfo
}

// generateServiceURL generates the service URL based on labels and configuration
func (c *Client) generateServiceURL(service ServiceInfo) string {
	if service.Project == "" || service.Service == "" {
		return ""
	}

	// Handle special services with known URL patterns
	switch service.Service {
	case "phpmyadmin":
		return fmt.Sprintf("http://pma.%s.localhost", service.Project)
	case "mailpit":
		return fmt.Sprintf("http://mail.%s.localhost", service.Project)
	case "app":
		return fmt.Sprintf("http://%s.localhost", service.Project)
	case "traefik":
		if service.Project == "phpier" {
			// Check if port 8080 is exposed for dashboard
			for _, port := range service.Ports {
				if port.ContainerPort == "8080" {
					return "http://localhost:8080"
				}
			}
		}
		return ""
	default:
		return ""
	}
}

// formatUptime formats a duration into a human-readable uptime string
func formatUptime(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours == 0 {
		return fmt.Sprintf("%dm", minutes)
	}

	if hours < 24 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}

	days := hours / 24
	hours = hours % 24

	if days == 1 {
		return fmt.Sprintf("1d %dh %dm", hours, minutes)
	}

	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}

// GetGlobalServices returns information about global phpier services
func (c *Client) GetGlobalServices(ctx context.Context) ([]ServiceInfo, error) {
	filter := &ServicesFilter{
		Project: "phpier",
	}

	return c.GetPhpierServices(ctx, filter)
}

// GetProjectServices returns information about services for a specific project
func (c *Client) GetProjectServices(ctx context.Context, projectName string) ([]ServiceInfo, error) {
	filter := &ServicesFilter{
		Project: projectName,
	}

	return c.GetPhpierServices(ctx, filter)
}

// ServiceExists checks if a specific service exists and is running
func (c *Client) ServiceExists(ctx context.Context, projectName, serviceName string) (bool, error) {
	containerName := fmt.Sprintf("%s-%s-1", projectName, serviceName)

	output, err := c.RunCommandOutput("docker", "ps", "-q", "--filter", fmt.Sprintf("name=%s", containerName))
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(output) != "", nil
}

// GetServiceStatus returns the status of a specific service
func (c *Client) GetServiceStatus(ctx context.Context, projectName, serviceName string) (string, error) {
	containerName := fmt.Sprintf("%s-%s-1", projectName, serviceName)

	output, err := c.RunCommandOutput("docker", "inspect", "--format", "{{.State.Status}}", containerName)
	if err != nil {
		return "not found", nil
	}

	return strings.TrimSpace(output), nil
}
