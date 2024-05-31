package middleware

import (
	"github.com/SoleaEnergy/forwardAuth/internal/policy"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestUserVerificationHandler(t *testing.T) func() {
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

func TestUserVerificationHandler_Success(t *testing.T) {
	teardown := setupTestUserVerificationHandler(t)
	defer teardown()

	handler := UserVerificationHandler(http.HandlerFunc(SuccessHandler))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"preferred_username": "testuser@example.com",
	})
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	req = mockPolicyFromRequest(req, &policy.Policy{
		Users: []string{"testuser"},
	})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserVerificationHandler_Unauthorized(t *testing.T) {
	teardown := setupTestUserVerificationHandler(t)
	defer teardown()

	handler := UserVerificationHandler(http.HandlerFunc(SuccessHandler))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"preferred_username": "unauthorizeduser@example.com",
	})
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	req = mockPolicyFromRequest(req, &policy.Policy{
		Users: []string{"testuser"},
	})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	actualVal := GetUserPolicyStatusFromContext(req)

	assert.Equal(t, false, actualVal)
}
