package updater

import (
	"testing"
)

func TestNewUpdater(t *testing.T) {
	version := "v1.0.0"
	verbose := true

	updater := NewUpdater(version, verbose)

	if updater.currentVer != version {
		t.Errorf("NewUpdater() currentVer = %v, want %v", updater.currentVer, version)
	}

	if updater.verbose != verbose {
		t.Errorf("NewUpdater() verbose = %v, want %v", updater.verbose, verbose)
	}
}

func TestExtractVersionFromTitle(t *testing.T) {
	updater := &Updater{}

	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "simple version with v prefix",
			title:    "v1.2.3",
			expected: "v1.2.3",
		},
		{
			name:     "simple version without v prefix",
			title:    "1.2.3",
			expected: "v1.2.3",
		},
		{
			name:     "release title with version",
			title:    "Release v2.1.0",
			expected: "v2.1.0",
		},
		{
			name:     "version with prerelease tag",
			title:    "v1.0.0-beta",
			expected: "v1.0.0-beta",
		},
		{
			name:     "complex release title",
			title:    "PHPier v1.5.2 - Bug fixes and improvements",
			expected: "v1.5.2",
		},
		{
			name:     "no version found",
			title:    "Some random title",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := updater.extractVersionFromTitle(tt.title)
			if result != tt.expected {
				t.Errorf("extractVersionFromTitle(%q) = %v, want %v", tt.title, result, tt.expected)
			}
		})
	}
}

func TestNeedsUpdate(t *testing.T) {
	updater := &Updater{}

	tests := []struct {
		name     string
		current  string
		latest   string
		expected bool
	}{
		{
			name:     "same version",
			current:  "v1.0.0",
			latest:   "v1.0.0",
			expected: false,
		},
		{
			name:     "newer version available",
			current:  "v1.0.0",
			latest:   "v1.0.1",
			expected: true,
		},
		{
			name:     "mixed v prefix",
			current:  "1.0.0",
			latest:   "v1.0.0",
			expected: false,
		},
		{
			name:     "current version unknown",
			current:  "unknown",
			latest:   "v1.0.0",
			expected: true,
		},
		{
			name:     "latest version empty",
			current:  "v1.0.0",
			latest:   "",
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
