package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEffectiveUser(t *testing.T) {
	tests := []struct {
		name     string
		shUser   string
		expected string
	}{
		{
			name:     "default user when shUser is empty",
			shUser:   "",
			expected: "www-data",
		},
		{
			name:     "custom user when shUser is set",
			shUser:   "root",
			expected: "root",
		},
		{
			name:     "custom user with non-standard name",
			shUser:   "developer",
			expected: "developer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the global variable for the test
			originalShUser := shUser
			shUser = tt.shUser
			defer func() { shUser = originalShUser }()

			result := getEffectiveUser()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "substring at beginning",
			s:        "bash not found",
			substr:   "bash",
			expected: true,
		},
		{
			name:     "substring at end",
			s:        "executable not found: bash",
			substr:   "bash",
			expected: true,
		},
		{
			name:     "substring in middle",
			s:        "error: bash failed to execute",
			substr:   "bash",
			expected: true,
		},
		{
			name:     "substring not found",
			s:        "shell not available",
			substr:   "bash",
			expected: false,
		},
		{
			name:     "empty substring",
			s:        "any string",
			substr:   "",
			expected: true,
		},
		{
			name:     "empty string",
			s:        "",
			substr:   "bash",
			expected: false,
		},
		{
			name:     "exact match",
			s:        "bash",
			substr:   "bash",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsAt(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "substring found",
			s:        "hello world bash test",
			substr:   "bash",
			expected: true,
		},
		{
			name:     "substring not found",
			s:        "hello world",
			substr:   "bash",
			expected: false,
		},
		{
			name:     "substring longer than string",
			s:        "hi",
			substr:   "hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsAt(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
