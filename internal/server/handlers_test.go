package server

import (
	"context"
	"github.com/lestrrat-go/jwx/jwk"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/SoleaEnergy/forwardAuth/internal/policy"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

const policyKey = "POLICY"

func setup(t *testing.T) func() {
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

func TestSuccessHandler(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()

	SuccessHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserVerificationHandler_Success(t *testing.T) {
	teardown := setup(t)
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
	teardown := setup(t)
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

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoleVerificationHandler_Success(t *testing.T) {
	teardown := setup(t)
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
	teardown := setup(t)
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

func TestJwtValidationHandler_Success(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	handler := JwtValidationHandler(http.HandlerFunc(SuccessHandler))

	// Mock a valid JWKS
	jwks = jwk.NewSet()
	key, _ := jwk.New([]byte("secret"))
	key.Set("kid", "test-key-id")
	jwks.Add(key)

	// Create a token with the "kid" in the header
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"roles": []interface{}{"admin"},
	})
	token.Header["kid"] = "test-key-id"
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJwtValidationHandler_Unauthorized_InvalidKid(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	handler := JwtValidationHandler(http.HandlerFunc(SuccessHandler))

	// Mock a valid JWKS
	jwks = jwk.NewSet()
	key, _ := jwk.New([]byte("secret"))
	key.Set("kid", "test-key-id")
	jwks.Add(key)

	// Create a token with the "kid" in the header
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"roles": []interface{}{"admin"},
	})
	token.Header["kid"] = "bad-kid"
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJwtValidationHandler_Unauthorized_MissingAlg(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	handler := JwtValidationHandler(http.HandlerFunc(SuccessHandler))

	// Mock a valid JWKS
	jwks = jwk.NewSet()
	key, _ := jwk.New([]byte("secret"))
	key.Set("kid", "test-key-id")
	jwks.Add(key)

	// Create a token with the "kid" in the header
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"roles": []interface{}{"admin"},
	})
	token.Header["kid"] = "test-key-id"
	token.Header["alg"] = nil
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
