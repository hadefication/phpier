package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display version information for phpier including build details.

This command shows:
- Version number
- Git commit hash  
- Build date
- Go version used to build
- Platform information`,
	Run: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("PHPier %s\n", buildVersion)
	fmt.Printf("Commit: %s\n", buildCommit)
	fmt.Printf("Built: %s\n", buildDate)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
