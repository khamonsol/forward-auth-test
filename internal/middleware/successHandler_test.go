package middleware

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestSuccessHandler(t *testing.T) func() {
	t.Helper()
	// Set environment variables for testing
	os.Setenv("ACCESS_TOKEN_HEADER", "Authorization")
	os.Setenv("JWKS_URL", "https://example.com/.well-known/jwks.json")
	os.Setenv("VALID_ALGS", "RS256,RS384,RS512")

	return func() {
		// Unset environment variables after the test
		os.Unsetenv("ACCESS_TOKEN_HEADER")
		os.Unsetenv("JWKS_URL")
	}
}

func TestSuccessHandler_StatusOK(t *testing.T) {
	teardown := setupTestSuccessHandler(t)
	defer teardown()

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	//Mock a passing Role check.
	newCtx := context.WithValue(req.Context(), rolePolicyCheckKey, true)
	w := httptest.NewRecorder()

	SuccessHandler(w, req.WithContext(newCtx))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSuccessHandler_StatusUnauthorized(t *testing.T) {
	teardown := setupTestSuccessHandler(t)
	defer teardown()

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	//Mock a passing Role check.
	newCtx := context.WithValue(req.Context(), rolePolicyCheckKey, false)
	newCtx = context.WithValue(newCtx, userPolicyCheckKey, false)
	w := httptest.NewRecorder()

	SuccessHandler(w, req.WithContext(newCtx))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
