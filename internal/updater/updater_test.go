package updater

import (
	"testing"
)

func TestNeedsUpdate(t *testing.T) {
	updater := &Updater{
		currentVer: "1.0.0",
	}

	tests := []struct {
		name     string
		current  string
		latest   string
		expected bool
	}{
		{
			name:     "same version",
			current:  "1.0.0",
			latest:   "1.0.0",
			expected: false,
		},
		{
			name:     "same version with v prefix",
			current:  "v1.0.0",
			latest:   "v1.0.0",
			expected: false,
		},
		{
			name:     "mixed v prefix",
			current:  "1.0.0",
			latest:   "v1.0.0",
			expected: false,
		},
		{
			name:     "newer version available",
			current:  "1.0.0",
			latest:   "1.0.1",
			expected: true,
		},
		{
			name:     "empty current version",
			current:  "",
			latest:   "1.0.0",
			expected: true,
		},
		{
			name:     "unknown current version",
			current:  "unknown",
			latest:   "1.0.0",
			expected: true,
		},
		{
			name:     "empty latest version",
			current:  "1.0.0",
			latest:   "",
			expected: false,
		},
		{
			name:     "unknown latest version",
			current:  "1.0.0",
			latest:   "unknown",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := updater.needsUpdate(tt.current, tt.latest)
			if result != tt.expected {
				t.Errorf("needsUpdate(%q, %q) = %v, want %v", tt.current, tt.latest, result, tt.expected)
			}
		})
	}
}

func TestIsValidVersionFormat(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "valid version with v prefix",
			version:  "v1.0.0",
			expected: true,
		},
		{
			name:     "valid version without v prefix",
			version:  "1.0.0",
			expected: true,
		},
		{
			name:     "valid version with patch",
			version:  "v1.2.3",
			expected: true,
		},
		{
			name:     "empty version",
			version:  "",
			expected: false,
		},
		{
			name:     "single character v",
			version:  "v",
			expected: false,
		},
		{
			name:     "invalid start character",
			version:  "x1.0.0",
			expected: false,
		},
		{
			name:     "valid major version only",
			version:  "v1",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidVersionFormat(tt.version)
			if result != tt.expected {
				t.Errorf("IsValidVersionFormat(%q) = %v, want %v", tt.version, result, tt.expected)
			}
		})
	}
}
