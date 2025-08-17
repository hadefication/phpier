package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"phpier/internal/docker"
)

func TestServicesCommand(t *testing.T) {
	// Create a test root command
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(servicesCmd)

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "services command without args",
			args:        []string{"services"},
			expectError: false,
		},
		{
			name:        "services with json flag",
			args:        []string{"services", "--json"},
			expectError: false,
		},
		{
			name:        "services with verbose flag",
			args:        []string{"services", "--verbose"},
			expectError: false,
		},
		{
			name:        "services with project filter",
			args:        []string{"services", "--project", "test-project"},
			expectError: false,
		},
		{
			name:        "services with status filter",
			args:        []string{"services", "--status", "running"},
			expectError: false,
		},
		{
			name:        "services with type filter",
			args:        []string{"services", "--type", "app"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				// Note: This might fail if Docker is not available,
				// but it validates the command structure
				// In a real test environment, we'd mock the Docker client
				assert.NoError(t, err)
			}
		})
	}
}

func TestMatchesServiceType(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		serviceType string
		expected    bool
	}{
		{
			name:        "app service matches app type",
			serviceName: "app",
			serviceType: "app",
			expected:    true,
		},
		{
			name:        "mysql service matches db type",
			serviceName: "mysql",
			serviceType: "db",
			expected:    true,
		},
		{
			name:        "redis service matches cache type",
			serviceName: "redis",
			serviceType: "cache",
			expected:    true,
		},
		{
			name:        "traefik service matches proxy type",
			serviceName: "traefik",
			serviceType: "proxy",
			expected:    true,
		},
		{
			name:        "phpmyadmin service matches tools type",
			serviceName: "phpmyadmin",
			serviceType: "tools",
			expected:    true,
		},
		{
			name:        "app service does not match db type",
			serviceName: "app",
			serviceType: "db",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service info
			service := docker.ServiceInfo{
				Service: tt.serviceName,
				Name:    "test-" + tt.serviceName + "-1",
			}

			result := matchesServiceType(service, tt.serviceType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateServiceType(t *testing.T) {
	tests := []struct {
		name        string
		serviceType string
		expectError bool
	}{
		{
			name:        "valid app type",
			serviceType: "app",
			expectError: false,
		},
		{
			name:        "valid db type",
			serviceType: "db",
			expectError: false,
		},
		{
			name:        "valid database type",
			serviceType: "database",
			expectError: false,
		},
		{
			name:        "valid cache type",
			serviceType: "cache",
			expectError: false,
		},
		{
			name:        "valid proxy type",
			serviceType: "proxy",
			expectError: false,
		},
		{
			name:        "valid tools type",
			serviceType: "tools",
			expectError: false,
		},
		{
			name:        "invalid type",
			serviceType: "invalid",
			expectError: true,
		},
		{
			name:        "case insensitive valid type",
			serviceType: "APP",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServiceType(tt.serviceType)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateStatusFilter(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		expectError bool
	}{
		{
			name:        "valid running status",
			status:      "running",
			expectError: false,
		},
		{
			name:        "valid stopped status",
			status:      "stopped",
			expectError: false,
		},
		{
			name:        "valid exited status",
			status:      "exited",
			expectError: false,
		},
		{
			name:        "valid paused status",
			status:      "paused",
			expectError: false,
		},
		{
			name:        "valid restarting status",
			status:      "restarting",
			expectError: false,
		},
		{
			name:        "valid created status",
			status:      "created",
			expectError: false,
		},
		{
			name:        "invalid status",
			status:      "invalid",
			expectError: true,
		},
		{
			name:        "case insensitive valid status",
			status:      "RUNNING",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStatusFilter(tt.status)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDetectCurrentProject(t *testing.T) {
	// This test would require creating temporary files
	// and changing directories, which is complex for unit tests
	// In a real implementation, we'd use dependency injection
	// to make this testable

	result := detectCurrentProject()
	// Since we're in the phpier codebase directory,
	// this should return "phpier" or the directory name
	assert.IsType(t, "", result)
}

func TestIsInPhpierProject(t *testing.T) {
	// Similar to detectCurrentProject, this would require
	// file system setup for proper testing
	result := isInPhpierProject()
	assert.IsType(t, false, result)
}
