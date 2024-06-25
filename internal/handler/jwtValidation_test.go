package handler

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type JwtValidationTestSuite struct {
	suite.Suite
	signingKey []byte
}

func (suite *JwtValidationTestSuite) SetupTest() {
	// Define the signing key for tests
	suite.signingKey = []byte("secret")
}

func (suite *JwtValidationTestSuite) createToken(roles []string) string {
	claims := jwt.MapClaims{
		"roles": roles,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(suite.signingKey)
	return tokenString
}

func (suite *JwtValidationTestSuite) TestJwtValidation_ValidToken() {
	tokenString := suite.createToken([]string{"role1"})
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	ctx := context.WithValue(req.Context(), correlationIdKey, "test-corr-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := JwtValidation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed, ok := GetAllowed(r)
		suite.True(ok)
		suite.True(allowed)
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *JwtValidationTestSuite) TestJwtValidation_InvalidToken() {
	tokenString := "invalid-token"
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	ctx := context.WithValue(req.Context(), correlationIdKey, "test-corr-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := JwtValidation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed, ok := GetAllowed(r)
		suite.True(ok)
		suite.False(allowed)
		status, _ := GetErrorStatus(r)
		suite.Equal(http.StatusUnauthorized, status)
		w.WriteHeader(status)
	}))

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *JwtValidationTestSuite) TestJwtValidation_NoToken() {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	ctx := context.WithValue(req.Context(), correlationIdKey, "test-corr-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := JwtValidation(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed, ok := GetAllowed(r)
		suite.True(ok)
		suite.False(allowed)
		status, _ := GetErrorStatus(r)
		suite.Equal(http.StatusUnauthorized, status)
		w.WriteHeader(status)
	}))

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func TestJwtValidationTestSuite(t *testing.T) {
	suite.Run(t, new(JwtValidationTestSuite))
}
