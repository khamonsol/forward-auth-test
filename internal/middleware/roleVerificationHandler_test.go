package middleware

import (
	"context"
	"github.com/SoleaEnergy/forwardAuth/internal/policy"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const policyKey = "POLICY"

func setupRoleVerificationHandlerTest(t *testing.T) func() {
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

// mockPolicyFromRequest injects a policy into the request context.
func mockPolicyFromRequest(r *http.Request, policy *policy.Policy) *http.Request {
	ctx := context.WithValue(r.Context(), policyKey, policy)
	return r.WithContext(ctx)
}

func TestRoleVerificationHandler_Success(t *testing.T) {
	teardown := setupRoleVerificationHandlerTest(t)
	defer teardown()

	handler := RoleVerificationHandler(http.HandlerFunc(SuccessHandler))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"roles": []interface{}{"admin"},
	})
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	req = mockPolicyFromRequest(req, &policy.Policy{
		Roles: []string{"admin"},
	})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleVerificationHandler_Unauthorized(t *testing.T) {
	teardown := setupRoleVerificationHandlerTest(t)
	defer teardown()

	handler := RoleVerificationHandler(http.HandlerFunc(SuccessHandler))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"roles": []interface{}{"user"},
	})
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	req = mockPolicyFromRequest(req, &policy.Policy{
		Roles: []string{"admin"},
	})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoleVerificationHandler_NoRolesOnToken(t *testing.T) {
	teardown := setupRoleVerificationHandlerTest(t)
	defer teardown()

	handler := RoleVerificationHandler(http.HandlerFunc(SuccessHandler))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	req = mockPolicyFromRequest(req, &policy.Policy{
		Roles: []string{"admin"},
	})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoleVerificationHandler_NoRolesOnPolicy(t *testing.T) {
	teardown := setupRoleVerificationHandlerTest(t)
	defer teardown()

	handler := RoleVerificationHandler(http.HandlerFunc(SuccessHandler))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"roles": []interface{}{"user"},
	})
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	req = mockPolicyFromRequest(req, &policy.Policy{})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
