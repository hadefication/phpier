package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of error that occurred
type ErrorType string

const (
	// Configuration errors
	ErrorTypeInvalidConfig   ErrorType = "invalid_config"
	ErrorTypeConfigNotFound  ErrorType = "config_not_found"
	ErrorTypeConfigCorrupted ErrorType = "config_corrupted"

	// Docker errors
	ErrorTypeDockerNotFound      ErrorType = "docker_not_found"
	ErrorTypeDockerNotRunning    ErrorType = "docker_not_running"
	ErrorTypeDockerError         ErrorType = "docker_error"
	ErrorTypeDockerComposeError  ErrorType = "docker_compose_error"
	ErrorTypeContainerNotFound   ErrorType = "container_not_found"
	ErrorTypeContainerNotRunning ErrorType = "container_not_running"
	ErrorTypeBuildFailed         ErrorType = "build_failed"

	// File system errors
	ErrorTypeFileNotFound    ErrorType = "file_not_found"
	ErrorTypeFilePermission  ErrorType = "file_permission"
	ErrorTypeDirectoryExists ErrorType = "directory_exists"
	ErrorTypeTemplateError   ErrorType = "template_error"
	ErrorTypeFileSystemError ErrorType = "file_system_error"

	// Network errors
	ErrorTypePortConflict   ErrorType = "port_conflict"
	ErrorTypeNetworkTimeout ErrorType = "network_timeout"

	// Validation errors
	ErrorTypeInvalidPHPVersion    ErrorType = "invalid_php_version"
	ErrorTypeInvalidDatabaseType  ErrorType = "invalid_database_type"
	ErrorTypeRequiredFieldMissing ErrorType = "required_field_missing"
	ErrorTypeInvalidPortRange     ErrorType = "invalid_port_range"

	// Command execution errors
	ErrorTypeCommandFailed    ErrorType = "command_failed"
	ErrorTypeCommandNotFound  ErrorType = "command_not_found"
	ErrorTypeCommandTimeout   ErrorType = "command_timeout"
	ErrorTypeInvalidArguments ErrorType = "invalid_arguments"

	// User interaction errors
	ErrorTypeUserAborted ErrorType = "user_aborted"

	// Internal errors
	ErrorTypeInternal ErrorType = "internal_error"
	ErrorTypeUnknown  ErrorType = "unknown_error"
	ErrorTypeFatal    ErrorType = "fatal_error"
)

// PhpierError represents a structured error with context and suggestions
type PhpierError struct {
	Type        ErrorType
	Message     string
	Cause       error
	Context     map[string]interface{}
	Suggestions []string
}

// Error implements the error interface
func (e *PhpierError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error
func (e *PhpierError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target error type
func (e *PhpierError) Is(target error) bool {
	if t, ok := target.(*PhpierError); ok {
		return e.Type == t.Type
	}
	return false
}

// WithContext adds context information to the error
func (e *PhpierError) WithContext(key string, value interface{}) *PhpierError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithSuggestion adds a suggestion to the error
func (e *PhpierError) WithSuggestion(suggestion string) *PhpierError {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// GetContext returns a context value
func (e *PhpierError) GetContext(key string) (interface{}, bool) {
	if e.Context == nil {
		return nil, false
	}
	value, exists := e.Context[key]
	return value, exists
}

// FormatUserFriendly returns a user-friendly error message with suggestions
func (e *PhpierError) FormatUserFriendly() string {
	var parts []string

	// Add the main error message
	parts = append(parts, fmt.Sprintf("Error: %s", e.Message))

	// Add context if available
	if len(e.Context) > 0 {
		contextParts := []string{}
		for k, v := range e.Context {
			contextParts = append(contextParts, fmt.Sprintf("%s: %v", k, v))
		}
		if len(contextParts) > 0 {
			parts = append(parts, fmt.Sprintf("Context: %s", strings.Join(contextParts, ", ")))
		}
	}

	// Add suggestions if available
	if len(e.Suggestions) > 0 {
		parts = append(parts, "")
		parts = append(parts, "Suggestions:")
		for i, suggestion := range e.Suggestions {
			parts = append(parts, fmt.Sprintf("  %d. %s", i+1, suggestion))
		}
	}

	return strings.Join(parts, "\n")
}

// NewPhpierError creates a new PhpierError
func NewPhpierError(errorType ErrorType, message string) *PhpierError {
	return &PhpierError{
		Type:    errorType,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// WrapError wraps an existing error with a PhpierError
func WrapError(errorType ErrorType, message string, cause error) *PhpierError {
	return &PhpierError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}
