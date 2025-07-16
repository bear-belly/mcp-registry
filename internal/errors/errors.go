package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "validation"
	ErrorTypeAuthentication ErrorType = "authentication"
	ErrorTypeAuthorization  ErrorType = "authorization"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeDatabase       ErrorType = "database"
	ErrorTypeInternal       ErrorType = "internal"
	ErrorTypeBadRequest     ErrorType = "bad_request"
)

// AppError represents an application error with context
type AppError struct {
	Type        ErrorType   `json:"type"`
	Message     string      `json:"message"`
	UserMessage string      `json:"user_message,omitempty"`
	Details     interface{} `json:"details,omitempty"`
	StatusCode  int         `json:"-"`
	Cause       error       `json:"-"`
	Stack       string      `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// NewAppError creates a new application error with stack trace
func NewAppError(errorType ErrorType, message string, cause error) *AppError {
	// Capture stack trace
	buf := make([]byte, 1024)
	runtime.Stack(buf, false)

	err := &AppError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
		Stack:   string(buf),
	}

	// Set default status codes based on error type
	switch errorType {
	case ErrorTypeValidation, ErrorTypeBadRequest:
		err.StatusCode = http.StatusBadRequest
	case ErrorTypeAuthentication:
		err.StatusCode = http.StatusUnauthorized
	case ErrorTypeAuthorization:
		err.StatusCode = http.StatusForbidden
	case ErrorTypeNotFound:
		err.StatusCode = http.StatusNotFound
	case ErrorTypeDatabase, ErrorTypeInternal:
		err.StatusCode = http.StatusInternalServerError
	default:
		err.StatusCode = http.StatusInternalServerError
	}

	return err
}

// SetUserMessage sets a user-friendly message
func (e *AppError) SetUserMessage(message string) *AppError {
	e.UserMessage = message
	return e
}

// SetDetails adds additional details to the error
func (e *AppError) SetDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// SetStatusCode overrides the default status code
func (e *AppError) SetStatusCode(code int) *AppError {
	e.StatusCode = code
	return e
}

// ErrorResponse represents the JSON error response structure
type ErrorResponse struct {
	Error   string      `json:"error"`
	Type    ErrorType   `json:"type"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// WriteError writes an error response to the HTTP writer
func WriteError(w http.ResponseWriter, err error) {
	var statusCode int
	var response ErrorResponse

	// Check if it's our custom AppError
	if ae, ok := err.(*AppError); ok {
		statusCode = ae.StatusCode
		response = ErrorResponse{
			Error:   ae.Message,
			Type:    ae.Type,
			Message: ae.UserMessage,
			Details: ae.Details,
		}

		// Log internal errors with stack trace
		if ae.Type == ErrorTypeInternal || ae.Type == ErrorTypeDatabase {
			log.Printf("Internal error: %s\nStack: %s", ae.Error(), ae.Stack)
		}
	} else {
		// Handle standard errors
		statusCode = http.StatusInternalServerError
		response = ErrorResponse{
			Error:   "Internal server error",
			Type:    ErrorTypeInternal,
			Message: "An unexpected error occurred",
		}
		log.Printf("Unhandled error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Common error constructors for convenience

// NewValidationError creates a validation error
func NewValidationError(message string, details interface{}) *AppError {
	return NewAppError(ErrorTypeValidation, message, nil).
		SetUserMessage("Please check your input and try again").
		SetDetails(details)
}

// NewAuthenticationError creates an authentication error
func NewAuthenticationError(message string) *AppError {
	return NewAppError(ErrorTypeAuthentication, message, nil).
		SetUserMessage("Authentication required")
}

// NewAuthorizationError creates an authorization error
func NewAuthorizationError(message string) *AppError {
	return NewAppError(ErrorTypeAuthorization, message, nil).
		SetUserMessage("You don't have permission to perform this action")
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return NewAppError(ErrorTypeNotFound, fmt.Sprintf("%s not found", resource), nil).
		SetUserMessage("The requested resource was not found")
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeDatabase, message, cause).
		SetUserMessage("A database error occurred")
}

// NewInternalError creates an internal error
func NewInternalError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeInternal, message, cause).
		SetUserMessage("An internal error occurred")
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return NewAppError(ErrorTypeBadRequest, message, nil).
		SetUserMessage("Invalid request")
}
