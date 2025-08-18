package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindProjectByName(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "phpier-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test project directories
	project1Dir := filepath.Join(tempDir, "test-project-1")
	project2Dir := filepath.Join(tempDir, "test-project-2")

	require.NoError(t, os.MkdirAll(project1Dir, 0755))
	require.NoError(t, os.MkdirAll(project2Dir, 0755))

	// Create .phpier.yml files
	phpierYml := "version: '3.8'\nservices:\n  app:\n    image: nginx"
	require.NoError(t, os.WriteFile(filepath.Join(project1Dir, ".phpier.yml"), []byte(phpierYml), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(project2Dir, ".phpier.yml"), []byte(phpierYml), 0644))

	// Create mock projects
	mockProjects := []ProjectInfo{
		{Name: "test-project-1", Path: project1Dir},
		{Name: "test-project-2", Path: project2Dir},
	}

	tests := []struct {
		name        string
		projectName string
		projects    []ProjectInfo
		wantErr     bool
		wantPath    string
	}{
		{
			name:        "Find existing project",
			projectName: "test-project-1",
			projects:    mockProjects,
			wantErr:     false,
			wantPath:    project1Dir,
		},
		{
			name:        "Project not found",
			projectName: "non-existent-project",
			projects:    mockProjects,
			wantErr:     true,
		},
		{
			name:        "Multiple projects with same name",
			projectName: "duplicate",
			projects: []ProjectInfo{
				{Name: "duplicate", Path: "/path1"},
				{Name: "duplicate", Path: "/path2"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the core logic by directly calling findProjectInList
			result, err := findProjectInList(tt.projectName, tt.projects)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.projectName, result.Name)
				assert.Equal(t, tt.wantPath, result.Path)
			}
		})
	}
}

func TestExtractProjectInfo(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "phpier-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create .phpier.yml file
	phpierYml := "version: '3.8'\nservices:\n  app:\n    image: nginx"
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, ".phpier.yml"), []byte(phpierYml), 0644))

	result, err := extractProjectInfo(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, filepath.Base(tempDir), result.Name)
	assert.Equal(t, tempDir, result.Path)
}

func TestExtractProjectInfo_NoConfigFile(t *testing.T) {
	// Create a temporary directory without .phpier.yml
	tempDir, err := os.MkdirTemp("", "phpier-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	result, err := extractProjectInfo(tempDir)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no .phpier.yml file found")
}

func TestLoadProjectConfigFromPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "phpier-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	result, err := LoadProjectConfigFromPath(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, filepath.Base(tempDir), result.Name)
	assert.Equal(t, "8.3", result.PHP)
	assert.Equal(t, "lts", result.Node)
}

func TestScanForProjects(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "phpier-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create nested project directories
	project1Dir := filepath.Join(tempDir, "project1")
	project2Dir := filepath.Join(tempDir, "subdir", "project2")
	nonProjectDir := filepath.Join(tempDir, "not-a-project")

	require.NoError(t, os.MkdirAll(project1Dir, 0755))
	require.NoError(t, os.MkdirAll(project2Dir, 0755))
	require.NoError(t, os.MkdirAll(nonProjectDir, 0755))

	// Create .phpier.yml files in project directories
	phpierYml := "version: '3.8'\nservices:\n  app:\n    image: nginx"
	require.NoError(t, os.WriteFile(filepath.Join(project1Dir, ".phpier.yml"), []byte(phpierYml), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(project2Dir, ".phpier.yml"), []byte(phpierYml), 0644))

	var projects []ProjectInfo
	err = scanForProjects(tempDir, 0, 3, &projects)
	assert.NoError(t, err)
	assert.Len(t, projects, 2)

	// Check that we found both projects
	projectNames := make(map[string]bool)
	for _, project := range projects {
		projectNames[project.Name] = true
	}
	assert.True(t, projectNames["project1"])
	assert.True(t, projectNames["project2"])
	assert.False(t, projectNames["not-a-project"])
}
