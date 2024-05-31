package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHandleError tests the HandleError function to ensure it logs the error and sends
// the correct JSON response with a correlation ID.
func TestHandleError(t *testing.T) {
	// Create a ResponseRecorder to capture the HTTP response
	rr := httptest.NewRecorder()
	errMessage := "test error"
	statusCode := http.StatusInternalServerError

	// Call the HandleError function
	HandleError(rr, errMessage, statusCode)

	// Check that the status code was set correctly
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Check that the content type is set to application/json
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Decode the JSON response
	var response ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	// Check that the error message contains the expected text
	expectedErrorMsg := "Error occurred, please contact support. Reference correlation ID:"
	assert.True(t, strings.Contains(response.Error, expectedErrorMsg))

	// Check that the correlation ID is a valid UUID without hyphens
	assert.Len(t, response.CorrelationID, 32)
	assert.NotContains(t, response.CorrelationID, "-")
}
