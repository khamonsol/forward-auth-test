package middleware

import (
	"github.com/SoleaEnergy/forwardAuth/internal/policy"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupAuthRequestTransactionIdHandler(t *testing.T) func() {
	t.Helper()
	// Set environment variables for testing
	os.Setenv("ACCESS_TOKEN_HEADER", "Authorization")
	os.Setenv("JWKS_URL", "https://example.com/.well-known/jwks.json")
	os.Setenv("VALID_ALGS", "RS256,RS384,RS512")

	return func() {
		// Unset environment variables after the test
		os.Unsetenv("ACCESS_TOKEN_HEADER")
		os.Unsetenv("JWKS_URL")
		os.Unsetenv("VALID_ALGS")
	}
}

func TestAuthRequestTransactionIdHandler(t *testing.T) {
	teardown := setupAuthRequestTransactionIdHandler(t)
	defer teardown()
	// Variable to store the modified request context
	var modifiedReq *http.Request

	// Wrap the SuccessHandler to capture the request
	captureContextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		modifiedReq = r
		SuccessHandler(w, r)
	})

	handler := AuthRequestTransactionIdHandler(captureContextHandler)

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", "Bearer dummy-token")
	req = mockPolicyFromRequest(req, &policy.Policy{
		Users: []string{"testuser"},
	})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	actualVal := GetAuthRequestTxId(modifiedReq)
	assert.NotEmpty(t, actualVal, "Expected correlation ID to be set in the request context")
}
