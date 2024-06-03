package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestJwtValidationHandler(t *testing.T) func() {
	t.Helper()
	// Set environment variables for testing
	os.Setenv("ACCESS_TOKEN_HEADER", "Authorization")
	os.Setenv("JWKS_URL", "https://example.com/.well-known/jwks.json")
	os.Setenv("VALID_ALGS", "HS256")

	return func() {
		// Unset environment variables after the test
		os.Unsetenv("ACCESS_TOKEN_HEADER")
		os.Unsetenv("JWKS_URL")
	}
}

func TestJwtValidationHandler_Success(t *testing.T) {
	teardown := setupTestJwtValidationHandler(t)
	defer teardown()

	passedJwtValidation := false
	// If the jwtops validation hands off to the next handler, that means this test was a success.
	captureContextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passedJwtValidation = true
		SuccessHandler(w, r)
	})

	handler := JwtValidationHandler(captureContextHandler)

	// Mock a valid JWKS
	jwks = jwk.NewSet()
	key, _ := jwk.New([]byte("secret"))
	key.Set("kid", "test-key-id")
	jwks.Add(key)

	// Create a token with the "kid" in the header
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"roles": []interface{}{"user"},
	})
	token.Header["kid"] = "test-key-id"
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	// If the jwtops validation hands off to the next handler, that means the jwtops successfully validated.
	assert.True(t, passedJwtValidation)
}

func TestJwtValidationHandler_Unauthorized_InvalidKid(t *testing.T) {
	teardown := setupTestJwtValidationHandler(t)
	defer teardown()

	passedJwtValidation := false
	// If the jwtops validation hands off to the next handler, that means this test was a success.
	captureContextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passedJwtValidation = true
		SuccessHandler(w, r)
	})
	handler := JwtValidationHandler(captureContextHandler)

	// Mock a valid JWKS
	jwks = jwk.NewSet()
	key, _ := jwk.New([]byte("secret"))
	key.Set("kid", "test-key-id")
	jwks.Add(key)

	// Create a token with the "kid" in the header
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"roles": []interface{}{"admin"},
	})
	token.Header["kid"] = "bad-kid"
	tokenString, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	req.Header.Set("Authorization", tokenString)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.False(t, passedJwtValidation)
}

func TestJwtValidationHandler_Unauthorized_MissingAlg(t *testing.T) {
	teardown := setupTestJwtValidationHandler(t)
	defer teardown()
	passedJwtValidation := false
	// If the jwtops validation hands off to the next handler, that means this test was a success.
	captureContextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passedJwtValidation = true
		SuccessHandler(w, r)
	})
	handler := JwtValidationHandler(captureContextHandler)

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

	assert.False(t, passedJwtValidation)
}
