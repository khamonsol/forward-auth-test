// Package util contains utility functions and types for the application.
package util

import (
	"encoding/json"
	"fmt"
	"log/slog" // Assuming slog is your structured logger
	"net/http"
)

// ErrorResponse defines the structure of the JSON response body for error messages.
// It includes a descriptive error message and a unique correlation ID for tracing the error.
type ErrorResponse struct {
	Error         string `json:"error"`         // Error provides a brief description of the error encountered.
	CorrelationID string `json:"correlationId"` // CorrelationID uniquely identifies the error instance for troubleshooting.
}

// HandleError logs the error and sends a JSON response with a correlation ID.
// If an internal server error occurs, it logs the error and returns a 403 Forbidden status,
// using the correlation ID for internal tracking and future diagnostics.
func HandleError(w http.ResponseWriter, errMessage string, statusCode int, correlationID string) {
	// Log the error with slog along with the correlation ID
	slog.Error(errMessage, "correlationId", correlationID, "status", statusCode)

	// If an internal server error is detected, log it and switch to 403 Forbidden
	if statusCode == http.StatusInternalServerError {
		slog.Error("Internal server error occurred", "correlationId", correlationID, "originalStatus", statusCode)
		statusCode = http.StatusForbidden // Use 403 to prevent retries and clarify behavior with forward auth
	}

	// Set the content type to application/json for the error response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Create the error response with the correlation ID
	msg := fmt.Sprintf("Error occurred, please contact support. Reference correlation ID: %s", correlationID)
	response := ErrorResponse{
		Error:         msg,
		CorrelationID: correlationID,
	}

	// Encode the response as JSON and send it
	_ = json.NewEncoder(w).Encode(response)
}
