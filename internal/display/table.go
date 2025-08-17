package display

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"phpier/internal/docker"
)

// TableOptions represents options for table display
type TableOptions struct {
	ShowHeaders bool
	ColorOutput bool
	MaxWidth    int
}

// ServicesTableConfig represents configuration for services table display
type ServicesTableConfig struct {
	ShowPorts    bool
	ShowMounts   bool
	ShowNetworks bool
	ShowURLs     bool
	Verbose      bool
}

// RenderServicesTable renders a table of services
func RenderServicesTable(services []docker.ServiceInfo, config ServicesTableConfig, options TableOptions) string {
	if len(services) == 0 {
		return colorize("No phpier services found", color.FgYellow, options.ColorOutput)
	}

	// Group services by project
	globalServices := []docker.ServiceInfo{}
	projectServices := make(map[string][]docker.ServiceInfo)

	for _, service := range services {
		if service.Project == "phpier" || service.Project == "" {
			globalServices = append(globalServices, service)
		} else {
			if projectServices[service.Project] == nil {
				projectServices[service.Project] = []docker.ServiceInfo{}
			}
			projectServices[service.Project] = append(projectServices[service.Project], service)
		}
	}

	var output strings.Builder

	// Render header
	output.WriteString(colorize("PHPIER SERVICES STATUS\n\n", color.FgCyan, options.ColorOutput))

	// Render global services
	if len(globalServices) > 0 {
		output.WriteString(colorize("Global Services:\n", color.FgMagenta, options.ColorOutput))
		output.WriteString(renderServiceGroup(globalServices, config, options))
		output.WriteString("\n")
	}

	// Render project services
	for projectName, services := range projectServices {
		header := fmt.Sprintf("Project Services (%s):\n", projectName)
		output.WriteString(colorize(header, color.FgMagenta, options.ColorOutput))
		output.WriteString(renderServiceGroup(services, config, options))
		output.WriteString("\n")
	}

	return strings.TrimSpace(output.String())
}

// renderServiceGroup renders a group of services as a table
func renderServiceGroup(services []docker.ServiceInfo, config ServicesTableConfig, options TableOptions) string {
	if len(services) == 0 {
		return colorize("  No services running\n", color.FgYellow, options.ColorOutput)
	}

	// Calculate column widths
	colWidths := calculateColumnWidths(services, config)

	var output strings.Builder

	// Render header if requested
	if options.ShowHeaders {
		output.WriteString(renderTableHeader(colWidths, config, options))
	}

	// Render each service
	for _, service := range services {
		output.WriteString(renderServiceRow(service, colWidths, config, options))
	}

	return output.String()
}

// calculateColumnWidths calculates the optimal column widths
func calculateColumnWidths(services []docker.ServiceInfo, config ServicesTableConfig) map[string]int {
	widths := map[string]int{
		"name":   8, // "NAME" header
		"status": 7, // "STATUS" header
		"uptime": 7, // "UPTIME" header
		"ports":  6, // "PORTS" header
		"image":  5, // "IMAGE" header
		"urls":   4, // "URLS" header
	}

	for _, service := range services {
		// Name column
		if len(service.Name) > widths["name"] {
			widths["name"] = len(service.Name)
		}

		// Status column
		if len(service.Status) > widths["status"] {
			widths["status"] = len(service.Status)
		}

		// Uptime column
		if len(service.Uptime) > widths["uptime"] {
			widths["uptime"] = len(service.Uptime)
		}

		// Ports column
		if config.ShowPorts {
			portStr := formatPorts(service.Ports)
			if len(portStr) > widths["ports"] {
				widths["ports"] = len(portStr)
			}
		}

		// Image column
		imageShort := formatImageName(service.Image)
		if len(imageShort) > widths["image"] {
			widths["image"] = len(imageShort)
		}

		// URLs column
		if config.ShowURLs && service.ServiceURL != "" {
			if len(service.ServiceURL) > widths["urls"] {
				widths["urls"] = len(service.ServiceURL)
			}
		}
	}

	// Set reasonable maximums
	if widths["name"] > 20 {
		widths["name"] = 20
	}
	if widths["image"] > 25 {
		widths["image"] = 25
	}
	if widths["urls"] > 40 {
		widths["urls"] = 40
	}

	return widths
}

// renderTableHeader renders the table header
func renderTableHeader(colWidths map[string]int, config ServicesTableConfig, options TableOptions) string {
	headers := []string{
		padRight("NAME", colWidths["name"]),
		padRight("STATUS", colWidths["status"]),
		padRight("UPTIME", colWidths["uptime"]),
	}

	if config.ShowPorts {
		headers = append(headers, padRight("PORTS", colWidths["ports"]))
	}

	if config.ShowURLs {
		headers = append(headers, padRight("URLS", colWidths["urls"]))
	}

	headers = append(headers, "IMAGE")

	headerLine := strings.Join(headers, "   ")
	return colorize(headerLine+"\n", color.FgWhite, options.ColorOutput)
}

// renderServiceRow renders a single service row
func renderServiceRow(service docker.ServiceInfo, colWidths map[string]int, config ServicesTableConfig, options TableOptions) string {
	// Format name (truncate if too long)
	name := truncateString(service.Name, colWidths["name"])

	// Format status with color
	status := service.Status
	if status == "" {
		status = "unknown"
	}

	// Format uptime
	uptime := service.Uptime
	if uptime == "" && service.Status == "running" {
		uptime = "unknown"
	} else if service.Status != "running" {
		uptime = "-"
	}

	// Build row components
	columns := []string{
		padRight(name, colWidths["name"]),
		padRight(formatStatus(status, options.ColorOutput), colWidths["status"]),
		padRight(uptime, colWidths["uptime"]),
	}

	if config.ShowPorts {
		ports := formatPorts(service.Ports)
		if ports == "" {
			ports = "-"
		}
		columns = append(columns, padRight(ports, colWidths["ports"]))
	}

	if config.ShowURLs {
		url := service.ServiceURL
		if url == "" {
			url = "-"
		}
		columns = append(columns, padRight(truncateString(url, colWidths["urls"]), colWidths["urls"]))
	}

	// Image column (last, no padding needed)
	image := formatImageName(service.Image)
	columns = append(columns, image)

	return strings.Join(columns, "   ") + "\n"
}

// formatStatus formats the status with appropriate colors
func formatStatus(status string, colorOutput bool) string {
	if !colorOutput {
		return status
	}

	switch strings.ToLower(status) {
	case "running":
		return color.GreenString(status)
	case "exited":
		return color.RedString(status)
	case "paused":
		return color.YellowString(status)
	case "restarting":
		return color.YellowString(status)
	case "created":
		return color.CyanString(status)
	default:
		return color.WhiteString(status)
	}
}

// formatPorts formats port mappings for display
func formatPorts(ports []docker.PortMapping) string {
	if len(ports) == 0 {
		return ""
	}

	var portStrs []string
	for _, port := range ports {
		if port.HostPort != "" {
			portStrs = append(portStrs, fmt.Sprintf("%s:%s", port.HostPort, port.ContainerPort))
		} else {
			portStrs = append(portStrs, port.ContainerPort)
		}
	}

	// Limit to first 3 ports to keep table readable
	if len(portStrs) > 3 {
		portStrs = append(portStrs[:3], "...")
	}

	return strings.Join(portStrs, ", ")
}

// formatImageName formats image names for display
func formatImageName(image string) string {
	if image == "" {
		return "-"
	}

	// Remove digest if present
	if strings.Contains(image, "@sha256:") {
		parts := strings.Split(image, "@")
		image = parts[0]
	}

	// For long image names, show only the relevant part
	if strings.Contains(image, "/") {
		parts := strings.Split(image, "/")
		if len(parts) > 2 {
			// Keep last two parts (e.g., "phpier/php:8.3")
			return strings.Join(parts[len(parts)-2:], "/")
		}
	}

	return image
}

// padRight pads a string to the right with spaces
func padRight(str string, length int) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(" ", length-len(str))
}

// truncateString truncates a string to the specified length
func truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	if maxLen <= 3 {
		return str[:maxLen]
	}
	return str[:maxLen-3] + "..."
}

// colorize applies color to text if color output is enabled
func colorize(text string, colorAttr color.Attribute, colorOutput bool) string {
	if !colorOutput {
		return text
	}
	return color.New(colorAttr).Sprint(text)
}

// RenderVerboseServicesInfo renders detailed service information
func RenderVerboseServicesInfo(services []docker.ServiceInfo, options TableOptions) string {
	if len(services) == 0 {
		return colorize("No phpier services found", color.FgYellow, options.ColorOutput)
	}

	var output strings.Builder

	for i, service := range services {
		if i > 0 {
			output.WriteString("\n")
		}
		output.WriteString(renderVerboseServiceInfo(service, options))
	}

	return output.String()
}

// renderVerboseServiceInfo renders detailed information for a single service
func renderVerboseServiceInfo(service docker.ServiceInfo, options TableOptions) string {
	var output strings.Builder

	// Service header
	header := fmt.Sprintf("=== %s ===", service.Name)
	output.WriteString(colorize(header+"\n", color.FgCyan, options.ColorOutput))

	// Basic information
	output.WriteString(fmt.Sprintf("Status: %s\n", formatStatus(service.Status, options.ColorOutput)))
	output.WriteString(fmt.Sprintf("Image: %s\n", service.Image))

	if service.Project != "" {
		output.WriteString(fmt.Sprintf("Project: %s\n", service.Project))
	}
	if service.Service != "" {
		output.WriteString(fmt.Sprintf("Service: %s\n", service.Service))
	}
	if service.Uptime != "" {
		output.WriteString(fmt.Sprintf("Uptime: %s\n", service.Uptime))
	}
	if service.ServiceURL != "" {
		output.WriteString(fmt.Sprintf("URL: %s\n", service.ServiceURL))
	}

	// Health status
	if service.Health != "" {
		output.WriteString(fmt.Sprintf("Health: %s\n", service.Health))
	}

	// Port mappings
	if len(service.Ports) > 0 {
		output.WriteString("Ports:\n")
		for _, port := range service.Ports {
			if port.HostPort != "" {
				output.WriteString(fmt.Sprintf("  %s:%s/%s\n", port.HostPort, port.ContainerPort, port.Protocol))
			} else {
				output.WriteString(fmt.Sprintf("  %s/%s\n", port.ContainerPort, port.Protocol))
			}
		}
	}

	// Networks
	if len(service.Networks) > 0 {
		output.WriteString(fmt.Sprintf("Networks: %s\n", strings.Join(service.Networks, ", ")))
	}

	// Mounts
	if len(service.Mounts) > 0 {
		output.WriteString("Mounts:\n")
		for _, mount := range service.Mounts {
			output.WriteString(fmt.Sprintf("  %s -> %s (%s)\n", mount.Source, mount.Destination, mount.Type))
		}
	}

	return output.String()
}
