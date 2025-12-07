package utils

import (
	"fmt"
	"runtime"
	"strings"
)

// Error represents a custom error with additional context
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	File    string `json:"file,omitempty"`
	Line    int    `json:"line,omitempty"`
	Func    string `json:"func,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// WithDetails adds additional details to the error
func (e *Error) WithDetails(details string) *Error {
	e.Details = details
	return e
}

// WithLocation adds file location information to the error
func (e *Error) WithLocation() *Error {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		e.File = file
		e.Line = line
		if fn := runtime.FuncForPC(pc); fn != nil {
			e.Func = fn.Name()
		}
	}
	return e
}

// NewError creates a new error with the given code and message
func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, code, message string) *Error {
	if err == nil {
		return nil
	}

	details := err.Error()
	if customErr, ok := err.(*Error); ok {
		details = customErr.Details
		if details == "" {
			details = customErr.Message
		}
	}

	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// IsError checks if an error has a specific code
func IsError(err error, code string) bool {
	if customErr, ok := err.(*Error); ok {
		return customErr.Code == code
	}
	return false
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) string {
	if customErr, ok := err.(*Error); ok {
		return customErr.Code
	}
	return "UNKNOWN"
}

// Common error codes
const (
	ErrCodeValidation        = "VALIDATION_ERROR"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeInternal          = "INTERNAL_ERROR"
	ErrCodeTimeout           = "TIMEOUT"
	ErrCodeRateLimit         = "RATE_LIMIT"
	ErrCodeInvalidInput      = "INVALID_INPUT"
	ErrCodeResourceExhausted = "RESOURCE_EXHAUSTED"
)

// Predefined errors
var (
	ErrValidation        = NewError(ErrCodeValidation, "Validation failed")
	ErrNotFound          = NewError(ErrCodeNotFound, "Resource not found")
	ErrUnauthorized      = NewError(ErrCodeUnauthorized, "Unauthorized access")
	ErrForbidden         = NewError(ErrCodeForbidden, "Access forbidden")
	ErrConflict          = NewError(ErrCodeConflict, "Resource conflict")
	ErrInternal          = NewError(ErrCodeInternal, "Internal server error")
	ErrTimeout           = NewError(ErrCodeTimeout, "Operation timeout")
	ErrRateLimit         = NewError(ErrCodeRateLimit, "Rate limit exceeded")
	ErrInvalidInput      = NewError(ErrCodeInvalidInput, "Invalid input provided")
	ErrResourceExhausted = NewError(ErrCodeResourceExhausted, "Resource exhausted")
)

// ErrorCollector collects multiple errors
type ErrorCollector struct {
	errors []error
}

// NewErrorCollector creates a new error collector
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make([]error, 0),
	}
}

// Add adds an error to the collector
func (ec *ErrorCollector) Add(err error) {
	if err != nil {
		ec.errors = append(ec.errors, err)
	}
}

// HasErrors returns true if there are any errors
func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.errors) > 0
}

// Errors returns all collected errors
func (ec *ErrorCollector) Errors() []error {
	return ec.errors
}

// Error returns a combined error message
func (ec *ErrorCollector) Error() string {
	if len(ec.errors) == 0 {
		return ""
	}

	var messages []string
	for _, err := range ec.errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// ToError returns the collector as a single error
func (ec *ErrorCollector) ToError() error {
	if len(ec.errors) == 0 {
		return nil
	}
	return NewError(ErrCodeValidation, "Multiple errors occurred").WithDetails(ec.Error())
}
