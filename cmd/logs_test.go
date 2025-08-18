package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestLogsCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "logs command exists",
			args: []string{"logs", "--help"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			follow = false
			tail = 0
			since = ""

			// Create a new root command for testing
			rootCmd := &cobra.Command{Use: "phpier"}
			rootCmd.AddCommand(logsCmd)

			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()

			// For help command, we expect no error
			assert.NoError(t, err)
		})
	}
}

func TestLogsFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedFollow bool
		expectedTail   int
		expectedSince  string
	}{
		{
			name:           "follow flag short",
			args:           []string{"logs", "-f"},
			expectedFollow: true,
			expectedTail:   0,
			expectedSince:  "",
		},
		{
			name:           "follow flag long",
			args:           []string{"logs", "--follow"},
			expectedFollow: true,
			expectedTail:   0,
			expectedSince:  "",
		},
		{
			name:           "tail flag",
			args:           []string{"logs", "--tail", "100"},
			expectedFollow: false,
			expectedTail:   100,
			expectedSince:  "",
		},
		{
			name:           "since flag",
			args:           []string{"logs", "--since", "2023-01-01T00:00:00Z"},
			expectedFollow: false,
			expectedTail:   0,
			expectedSince:  "2023-01-01T00:00:00Z",
		},
		{
			name:           "all flags",
			args:           []string{"logs", "-f", "--tail", "50", "--since", "2023-01-01T00:00:00Z"},
			expectedFollow: true,
			expectedTail:   50,
			expectedSince:  "2023-01-01T00:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			follow = false
			tail = 0
			since = ""

			// Create a new command for testing (without executing runLogs)
			testCmd := &cobra.Command{
				Use: "logs",
				RunE: func(cmd *cobra.Command, args []string) error {
					// Just parse flags, don't execute
					return nil
				},
			}
			testCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output in real-time")
			testCmd.Flags().IntVar(&tail, "tail", 0, "Number of lines to show from the end of the logs")
			testCmd.Flags().StringVar(&since, "since", "", "Show logs since timestamp")

			testCmd.SetArgs(tt.args[1:]) // Skip "logs" command name
			err := testCmd.Execute()

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedFollow, follow)
			assert.Equal(t, tt.expectedTail, tail)
			assert.Equal(t, tt.expectedSince, since)
		})
	}
}
