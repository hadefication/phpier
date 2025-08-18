package cmd

import (
	"fmt"
	"phpier/internal/config"
	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	showDocker     bool
	showFilesystem bool
	showAll        bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all discovered phpier projects",
	Long: `List all phpier projects discovered through Docker images and filesystem scanning.

This command helps you see which projects are available for use with global commands.

Examples:
  phpier list                    # List all projects from Docker and filesystem
  phpier list --docker           # List only projects found in Docker images
  phpier list --filesystem       # List only projects found in filesystem
  phpier list --all              # Show detailed information for all projects`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Flags
	listCmd.Flags().BoolVar(&showDocker, "docker", false, "Show only projects discovered from Docker images")
	listCmd.Flags().BoolVar(&showFilesystem, "filesystem", false, "Show only projects discovered from filesystem")
	listCmd.Flags().BoolVar(&showAll, "all", false, "Show detailed information including paths and status")
}

func runList(cmd *cobra.Command, args []string) error {
	var dockerProjects, filesystemProjects []config.ProjectInfo
	var err error

	// Determine what to show based on flags
	showDockerProjects := showDocker || (!showDocker && !showFilesystem)
	showFilesystemProjects := showFilesystem || (!showDocker && !showFilesystem)

	if showDockerProjects {
		logrus.Infof("🐳 Discovering projects from Docker images...")
		dockerProjects, err = config.DiscoverProjectsFromDocker()
		if err != nil {
			logrus.Warnf("Failed to discover Docker projects: %v", err)
		}
	}

	if showFilesystemProjects {
		logrus.Infof("📁 Discovering projects from filesystem...")
		filesystemProjects, err = config.DiscoverProjectsFromFilesystem()
		if err != nil {
			return errors.WrapError(errors.ErrorTypeProjectDiscoveryFailed, "Failed to discover filesystem projects", err)
		}
	}

	// Combine and deduplicate projects
	allProjects := make(map[string]config.ProjectInfo)
	
	// Add Docker projects
	for _, project := range dockerProjects {
		allProjects[project.Name] = project
	}
	
	// Add filesystem projects (Docker takes precedence)
	for _, project := range filesystemProjects {
		if existing, exists := allProjects[project.Name]; !exists || existing.Path == "" {
			allProjects[project.Name] = project
		}
	}

	if len(allProjects) == 0 {
		fmt.Println("No phpier projects found.")
		fmt.Println()
		fmt.Println("To create a new project:")
		fmt.Println("  phpier init 8.3 --project-name=my-project")
		return nil
	}

	// Display results
	if showAll {
		fmt.Printf("Found %d phpier project(s):\n\n", len(allProjects))
		
		if showDockerProjects && len(dockerProjects) > 0 {
			fmt.Println("🐳 Docker Projects:")
			for _, project := range dockerProjects {
				fmt.Printf("  %-20s %s\n", project.Name, formatProjectPath(project.Path))
			}
			fmt.Println()
		}
		
		if showFilesystemProjects && len(filesystemProjects) > 0 {
			fmt.Println("📁 Filesystem Projects:")
			for _, project := range filesystemProjects {
				fmt.Printf("  %-20s %s\n", project.Name, formatProjectPath(project.Path))
			}
			fmt.Println()
		}
	} else {
		fmt.Printf("Available phpier projects (%d):\n", len(allProjects))
		for name := range allProjects {
			fmt.Printf("  %s\n", name)
		}
		fmt.Println()
		fmt.Println("Use 'phpier list --all' for detailed information")
	}

	fmt.Println("Usage:")
	fmt.Println("  phpier up <project-name>      # Start a project")
	fmt.Println("  phpier down <project-name>    # Stop a project")

	return nil
}

func formatProjectPath(path string) string {
	if path == "" {
		return "(path unknown)"
	}
	return path
}