package errors

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// ExitCode represents exit codes for different error types
type ExitCode int

const (
	ExitCodeSuccess            ExitCode = 0
	ExitCodeGeneralError       ExitCode = 1
	ExitCodeDockerError        ExitCode = 2
	ExitCodeConfigurationError ExitCode = 3
	ExitCodeFileSystemError    ExitCode = 4
	ExitCodeValidationError    ExitCode = 5
	ExitCodeNetworkError       ExitCode = 6
	ExitCodeCommandError       ExitCode = 7
	ExitCodeInternalError      ExitCode = 8
)

// GetExitCode returns the appropriate exit code for an error type
func GetExitCode(errorType ErrorType) ExitCode {
	switch errorType {
	case ErrorTypeDockerNotFound, ErrorTypeDockerNotRunning, ErrorTypeDockerComposeError,
		ErrorTypeContainerNotFound, ErrorTypeContainerNotRunning, ErrorTypeBuildFailed:
		return ExitCodeDockerError

	case ErrorTypeInvalidConfig, ErrorTypeConfigNotFound, ErrorTypeConfigCorrupted:
		return ExitCodeConfigurationError

	case ErrorTypeFileNotFound, ErrorTypeFilePermission, ErrorTypeDirectoryExists, ErrorTypeTemplateError:
		return ExitCodeFileSystemError

	case ErrorTypeInvalidPHPVersion, ErrorTypeInvalidDatabaseType, ErrorTypeRequiredFieldMissing,
		ErrorTypeInvalidPortRange, ErrorTypePortConflict:
		return ExitCodeValidationError

	case ErrorTypeNetworkTimeout:
		return ExitCodeNetworkError

	case ErrorTypeCommandFailed, ErrorTypeCommandNotFound, ErrorTypeCommandTimeout, ErrorTypeInvalidArguments:
		return ExitCodeCommandError

	case ErrorTypeInternal:
		return ExitCodeInternalError

	default:
		return ExitCodeGeneralError
	}
}

// ErrorHandler handles error display and program termination
type ErrorHandler struct {
	verbose bool
	colored bool
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(verbose, colored bool) *ErrorHandler {
	return &ErrorHandler{
		verbose: verbose,
		colored: colored,
	}
}

// Handle processes and displays an error, then exits with appropriate code
func (h *ErrorHandler) Handle(err error) {
	if err == nil {
		return
	}

	var exitCode ExitCode = ExitCodeGeneralError

	// Check if it's a PhpierError
	if devErr, ok := err.(*PhpierError); ok {
		h.displayPhpierError(devErr)
		exitCode = GetExitCode(devErr.Type)
	} else {
		h.displayGenericError(err)
	}

	os.Exit(int(exitCode))
}

// HandleWithoutExit processes and displays an error without exiting
func (h *ErrorHandler) HandleWithoutExit(err error) ExitCode {
	if err == nil {
		return ExitCodeSuccess
	}

	var exitCode ExitCode = ExitCodeGeneralError

	// Check if it's a PhpierError
	if devErr, ok := err.(*PhpierError); ok {
		h.displayPhpierError(devErr)
		exitCode = GetExitCode(devErr.Type)
	} else {
		h.displayGenericError(err)
	}

	return exitCode
}

// displayPhpierError displays a PhpierError with formatting
func (h *ErrorHandler) displayPhpierError(err *PhpierError) {
	if h.colored {
		h.displayColoredPhpierError(err)
	} else {
		h.displayPlainPhpierError(err)
	}

	// Log detailed information if verbose mode is enabled
	if h.verbose {
		h.logVerboseErrorInfo(err)
	}
}

// displayColoredPhpierError displays a PhpierError with colors
func (h *ErrorHandler) displayColoredPhpierError(err *PhpierError) {
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)
	white := color.New(color.FgWhite)

	// Error header
	red.Fprint(os.Stderr, "✗ Error: ")
	fmt.Fprintf(os.Stderr, "%s\n", err.Message)

	// Context information
	if len(err.Context) > 0 {
		fmt.Fprint(os.Stderr, "\n")
		cyan.Fprint(os.Stderr, "Context:\n")
		for key, value := range err.Context {
			fmt.Fprintf(os.Stderr, "  %s: %v\n", key, value)
		}
	}

	// Suggestions
	if len(err.Suggestions) > 0 {
		fmt.Fprint(os.Stderr, "\n")
		yellow.Fprint(os.Stderr, "Suggestions:\n")
		for i, suggestion := range err.Suggestions {
			white.Fprintf(os.Stderr, "  %d. %s\n", i+1, suggestion)
		}
	}

	fmt.Fprint(os.Stderr, "\n")
}

// displayPlainPhpierError displays a PhpierError without colors
func (h *ErrorHandler) displayPlainPhpierError(err *PhpierError) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Message)

	// Context information
	if len(err.Context) > 0 {
		fmt.Fprint(os.Stderr, "\nContext:\n")
		for key, value := range err.Context {
			fmt.Fprintf(os.Stderr, "  %s: %v\n", key, value)
		}
	}

	// Suggestions
	if len(err.Suggestions) > 0 {
		fmt.Fprint(os.Stderr, "\nSuggestions:\n")
		for i, suggestion := range err.Suggestions {
			fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, suggestion)
		}
	}

	fmt.Fprint(os.Stderr, "\n")
}

// displayGenericError displays a generic error
func (h *ErrorHandler) displayGenericError(err error) {
	if h.colored {
		red := color.New(color.FgRed, color.Bold)
		red.Fprint(os.Stderr, "✗ Error: ")
	} else {
		fmt.Fprint(os.Stderr, "Error: ")
	}

	fmt.Fprintf(os.Stderr, "%s\n\n", err.Error())

	if h.verbose {
		logrus.WithError(err).Error("Detailed error information")
	}
}

// logVerboseErrorInfo logs detailed error information
func (h *ErrorHandler) logVerboseErrorInfo(err *PhpierError) {
	entry := logrus.WithFields(logrus.Fields{
		"error_type": string(err.Type),
		"context":    err.Context,
	})

	if err.Cause != nil {
		entry = entry.WithError(err.Cause)
	}

	entry.Debug("Detailed error information")
}

// IsPhpierError checks if an error is a PhpierError
func IsPhpierError(err error) bool {
	_, ok := err.(*PhpierError)
	return ok
}

// GetErrorType extracts the error type from a PhpierError
func GetErrorType(err error) ErrorType {
	if devErr, ok := err.(*PhpierError); ok {
		return devErr.Type
	}
	return ErrorTypeUnknown
}

// HasSuggestions checks if an error has suggestions
func HasSuggestions(err error) bool {
	if devErr, ok := err.(*PhpierError); ok {
		return len(devErr.Suggestions) > 0
	}
	return false
}
