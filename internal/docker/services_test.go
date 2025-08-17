package docker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "seconds only",
			duration: 30 * time.Second,
			expected: "30s",
		},
		{
			name:     "minutes only",
			duration: 5 * time.Minute,
			expected: "5m",
		},
		{
			name:     "hours and minutes",
			duration: 2*time.Hour + 30*time.Minute,
			expected: "2h 30m",
		},
		{
			name:     "one day",
			duration: 24*time.Hour + 2*time.Hour + 15*time.Minute,
			expected: "1d 2h 15m",
		},
		{
			name:     "multiple days",
			duration: 3*24*time.Hour + 5*time.Hour + 45*time.Minute,
			expected: "3d 5h 45m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatUptime(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPhpierContainer(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name          string
		containerName string
		expected      bool
	}{
		{
			name:          "phpier global traefik",
			containerName: "phpier-traefik-1",
			expected:      true,
		},
		{
			name:          "project app container",
			containerName: "myproject-app-1",
			expected:      true,
		},
		{
			name:          "project mysql container",
			containerName: "myproject-mysql-1",
			expected:      true,
		},
		{
			name:          "project redis container",
			containerName: "myproject-redis-1",
			expected:      true,
		},
		{
			name:          "project phpmyadmin container",
			containerName: "myproject-phpmyadmin-1",
			expected:      true,
		},
		{
			name:          "project mailpit container",
			containerName: "myproject-mailpit-1",
			expected:      true,
		},
		{
			name:          "global mysql container",
			containerName: "phpier-mysql-1",
			expected:      true,
		},
		{
			name:          "non-phpier container",
			containerName: "some-other-container",
			expected:      false,
		},
		{
			name:          "similar but not phpier container",
			containerName: "not-phpier-app-1",
			expected:      false,
		},
		{
			name:          "wordpress container",
			containerName: "wordpress-app-1",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.isPhpierContainer(tt.containerName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateServiceURL(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name     string
		service  ServiceInfo
		expected string
	}{
		{
			name: "app service",
			service: ServiceInfo{
				Project: "myproject",
				Service: "app",
			},
			expected: "http://myproject.localhost",
		},
		{
			name: "phpmyadmin service",
			service: ServiceInfo{
				Project: "myproject",
				Service: "phpmyadmin",
			},
			expected: "http://pma.myproject.localhost",
		},
		{
			name: "mailpit service",
			service: ServiceInfo{
				Project: "myproject",
				Service: "mailpit",
			},
			expected: "http://mail.myproject.localhost",
		},
		{
			name: "traefik service with dashboard port",
			service: ServiceInfo{
				Project: "phpier",
				Service: "traefik",
				Ports: []PortMapping{
					{ContainerPort: "8080", HostPort: "8080"},
				},
			},
			expected: "http://localhost:8080",
		},
		{
			name: "mysql service (no URL)",
			service: ServiceInfo{
				Project: "myproject",
				Service: "mysql",
			},
			expected: "",
		},
		{
			name: "service without project",
			service: ServiceInfo{
				Service: "app",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.generateServiceURL(tt.service)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractPorts(t *testing.T) {
	tests := []struct {
		name     string
		ports    map[string]interface{}
		expected []PortMapping
	}{
		{
			name: "single port mapping",
			ports: map[string]interface{}{
				"80/tcp": []interface{}{
					map[string]interface{}{
						"HostIp":   "0.0.0.0",
						"HostPort": "8080",
					},
				},
			},
			expected: []PortMapping{
				{
					ContainerPort: "80",
					HostPort:      "8080",
					Protocol:      "tcp",
				},
			},
		},
		{
			name: "multiple port mappings",
			ports: map[string]interface{}{
				"80/tcp": []interface{}{
					map[string]interface{}{
						"HostIp":   "0.0.0.0",
						"HostPort": "8080",
					},
				},
				"443/tcp": []interface{}{
					map[string]interface{}{
						"HostIp":   "0.0.0.0",
						"HostPort": "8443",
					},
				},
			},
			expected: []PortMapping{
				{
					ContainerPort: "80",
					HostPort:      "8080",
					Protocol:      "tcp",
				},
				{
					ContainerPort: "443",
					HostPort:      "8443",
					Protocol:      "tcp",
				},
			},
		},
		{
			name:     "no port mappings",
			ports:    map[string]interface{}{},
			expected: []PortMapping{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPorts(tt.ports)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestExtractMounts(t *testing.T) {
	tests := []struct {
		name     string
		mounts   []interface{}
		expected []MountInfo
	}{
		{
			name: "single mount",
			mounts: []interface{}{
				map[string]interface{}{
					"Source":      "/host/path",
					"Destination": "/container/path",
					"Type":        "bind",
					"Mode":        "rw",
				},
			},
			expected: []MountInfo{
				{
					Source:      "/host/path",
					Destination: "/container/path",
					Type:        "bind",
					Mode:        "rw",
				},
			},
		},
		{
			name: "multiple mounts",
			mounts: []interface{}{
				map[string]interface{}{
					"Source":      "/host/path1",
					"Destination": "/container/path1",
					"Type":        "bind",
					"Mode":        "rw",
				},
				map[string]interface{}{
					"Source":      "/host/path2",
					"Destination": "/container/path2",
					"Type":        "volume",
					"Mode":        "ro",
				},
			},
			expected: []MountInfo{
				{
					Source:      "/host/path1",
					Destination: "/container/path1",
					Type:        "bind",
					Mode:        "rw",
				},
				{
					Source:      "/host/path2",
					Destination: "/container/path2",
					Type:        "volume",
					Mode:        "ro",
				},
			},
		},
		{
			name:     "no mounts",
			mounts:   []interface{}{},
			expected: []MountInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractMounts(tt.mounts)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}
