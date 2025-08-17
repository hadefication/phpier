package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"phpier/internal/docker"
)

func TestFormatStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		colorOutput bool
		expected    string
	}{
		{
			name:        "running status without color",
			status:      "running",
			colorOutput: false,
			expected:    "running",
		},
		{
			name:        "exited status without color",
			status:      "exited",
			colorOutput: false,
			expected:    "exited",
		},
		{
			name:        "unknown status without color",
			status:      "unknown",
			colorOutput: false,
			expected:    "unknown",
		},
		// Color tests would require checking ANSI codes
		// which is complex, so we'll test the basic function
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatStatus(tt.status, tt.colorOutput)
			if !tt.colorOutput {
				assert.Equal(t, tt.expected, result)
			} else {
				// With color, result should contain the status text
				assert.Contains(t, result, tt.status)
			}
		})
	}
}

func TestFormatPorts(t *testing.T) {
	tests := []struct {
		name     string
		ports    []docker.PortMapping
		expected string
	}{
		{
			name:     "no ports",
			ports:    []docker.PortMapping{},
			expected: "",
		},
		{
			name: "single port",
			ports: []docker.PortMapping{
				{ContainerPort: "80", HostPort: "8080", Protocol: "tcp"},
			},
			expected: "8080:80",
		},
		{
			name: "multiple ports",
			ports: []docker.PortMapping{
				{ContainerPort: "80", HostPort: "8080", Protocol: "tcp"},
				{ContainerPort: "443", HostPort: "8443", Protocol: "tcp"},
			},
			expected: "8080:80, 8443:443",
		},
		{
			name: "port without host mapping",
			ports: []docker.PortMapping{
				{ContainerPort: "80", HostPort: "", Protocol: "tcp"},
			},
			expected: "80",
		},
		{
			name: "many ports (should truncate)",
			ports: []docker.PortMapping{
				{ContainerPort: "80", HostPort: "8080", Protocol: "tcp"},
				{ContainerPort: "443", HostPort: "8443", Protocol: "tcp"},
				{ContainerPort: "3306", HostPort: "3306", Protocol: "tcp"},
				{ContainerPort: "6379", HostPort: "6379", Protocol: "tcp"},
			},
			expected: "8080:80, 8443:443, 3306:3306, ...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPorts(tt.ports)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatImageName(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		expected string
	}{
		{
			name:     "empty image",
			image:    "",
			expected: "-",
		},
		{
			name:     "simple image name",
			image:    "nginx",
			expected: "nginx",
		},
		{
			name:     "image with tag",
			image:    "nginx:latest",
			expected: "nginx:latest",
		},
		{
			name:     "image with registry",
			image:    "docker.io/library/nginx:latest",
			expected: "library/nginx:latest",
		},
		{
			name:     "image with sha256 digest",
			image:    "nginx@sha256:1234567890abcdef",
			expected: "nginx",
		},
		{
			name:     "complex image name",
			image:    "registry.example.com/namespace/phpier/php:8.3",
			expected: "phpier/php:8.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatImageName(tt.image)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		maxLen   int
		expected string
	}{
		{
			name:     "string shorter than max",
			str:      "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "string equal to max",
			str:      "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "string longer than max",
			str:      "hello world",
			maxLen:   8,
			expected: "hello...",
		},
		{
			name:     "very short max length",
			str:      "hello",
			maxLen:   3,
			expected: "hel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.str, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		length   int
		expected string
	}{
		{
			name:     "string shorter than length",
			str:      "hello",
			length:   10,
			expected: "hello     ",
		},
		{
			name:     "string equal to length",
			str:      "hello",
			length:   5,
			expected: "hello",
		},
		{
			name:     "string longer than length",
			str:      "hello world",
			length:   5,
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padRight(tt.str, tt.length)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRenderServicesTable(t *testing.T) {
	services := []docker.ServiceInfo{
		{
			Name:    "test-app-1",
			Image:   "phpier/php:8.3",
			Status:  "running",
			Uptime:  "1h 30m",
			Project: "test",
			Service: "app",
			Ports: []docker.PortMapping{
				{ContainerPort: "80", HostPort: "8080", Protocol: "tcp"},
			},
			ServiceURL: "http://test.localhost",
		},
		{
			Name:    "test-mysql-1",
			Image:   "mysql:8.0",
			Status:  "running",
			Uptime:  "1h 30m",
			Project: "test",
			Service: "mysql",
			Ports: []docker.PortMapping{
				{ContainerPort: "3306", HostPort: "3306", Protocol: "tcp"},
			},
		},
	}

	config := ServicesTableConfig{
		ShowPorts: true,
		ShowURLs:  true,
	}

	options := TableOptions{
		ShowHeaders: true,
		ColorOutput: false,
	}

	result := RenderServicesTable(services, config, options)

	// Basic checks
	assert.Contains(t, result, "PHPIER SERVICES STATUS")
	assert.Contains(t, result, "Project Services (test)")
	assert.Contains(t, result, "test-app-1")
	assert.Contains(t, result, "test-mysql-1")
	assert.Contains(t, result, "running")
	assert.Contains(t, result, "1h 30m")
}

func TestRenderServicesTableEmpty(t *testing.T) {
	services := []docker.ServiceInfo{}
	config := ServicesTableConfig{}
	options := TableOptions{ColorOutput: false}

	result := RenderServicesTable(services, config, options)
	assert.Contains(t, result, "No phpier services found")
}

func TestRenderVerboseServicesInfo(t *testing.T) {
	services := []docker.ServiceInfo{
		{
			Name:       "test-app-1",
			Image:      "phpier/php:8.3",
			Status:     "running",
			Uptime:     "1h 30m",
			Project:    "test",
			Service:    "app",
			ServiceURL: "http://test.localhost",
			Health:     "healthy",
			Ports: []docker.PortMapping{
				{ContainerPort: "80", HostPort: "8080", Protocol: "tcp"},
			},
			Networks: []string{"test_default"},
			Mounts: []docker.MountInfo{
				{Source: "/host/path", Destination: "/var/www/html", Type: "bind"},
			},
		},
	}

	options := TableOptions{ColorOutput: false}
	result := RenderVerboseServicesInfo(services, options)

	// Verbose output checks
	assert.Contains(t, result, "=== test-app-1 ===")
	assert.Contains(t, result, "Status: running")
	assert.Contains(t, result, "Image: phpier/php:8.3")
	assert.Contains(t, result, "Project: test")
	assert.Contains(t, result, "Service: app")
	assert.Contains(t, result, "Uptime: 1h 30m")
	assert.Contains(t, result, "URL: http://test.localhost")
	assert.Contains(t, result, "Health: healthy")
	assert.Contains(t, result, "Ports:")
	assert.Contains(t, result, "8080:80/tcp")
	assert.Contains(t, result, "Networks: test_default")
	assert.Contains(t, result, "Mounts:")
	assert.Contains(t, result, "/host/path -> /var/www/html (bind)")
}

func TestCalculateColumnWidths(t *testing.T) {
	services := []docker.ServiceInfo{
		{
			Name:   "very-long-container-name-that-should-be-truncated",
			Status: "running",
			Uptime: "1h 30m",
			Image:  "some/very/long/image/name/that/exceeds/normal/length",
			Ports: []docker.PortMapping{
				{ContainerPort: "80", HostPort: "8080", Protocol: "tcp"},
			},
			ServiceURL: "http://very.long.domain.name.localhost",
		},
	}

	config := ServicesTableConfig{ShowPorts: true, ShowURLs: true}
	widths := calculateColumnWidths(services, config)

	// Check that widths are calculated
	assert.GreaterOrEqual(t, widths["name"], 8)
	assert.GreaterOrEqual(t, widths["status"], 7)
	assert.GreaterOrEqual(t, widths["uptime"], 7)
	assert.Greater(t, widths["ports"], 6)
	assert.Greater(t, widths["image"], 5)
	assert.Greater(t, widths["urls"], 4)

	// Check maximum limits are applied
	assert.LessOrEqual(t, widths["name"], 20)
	assert.LessOrEqual(t, widths["image"], 25)
	assert.LessOrEqual(t, widths["urls"], 40)
}
