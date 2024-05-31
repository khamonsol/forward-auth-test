package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("ACCESS_TOKEN_HEADER", "X-Access-Token")
	os.Setenv("JWKS_URL", "https://example.com/.well-known/jwks.json")

	// Ensure environment variables are unset after the test
	defer func() {
		os.Unsetenv("ACCESS_TOKEN_HEADER")
		os.Unsetenv("JWKS_URL")
	}()

	config := LoadConfig()

	assert.Equal(t, "X-Access-Token", config.AccessTokenHeader, "Expected ACCESS_TOKEN_HEADER to be set to 'X-Access-Token'")
	assert.Equal(t, "https://example.com/.well-known/jwks.json", config.JwksUrl, "Expected JWKS_URL to be set to 'https://example.com/.well-known/jwks.json'")
}

func TestLoadConfig_DefaultHeader(t *testing.T) {
	// Set only the required environment variable
	os.Setenv("JWKS_URL", "https://example.com/.well-known/jwks.json")

	// Ensure environment variables are unset after the test
	defer func() {
		os.Unsetenv("JWKS_URL")
	}()

	config := LoadConfig()

	assert.Equal(t, "Authorization", config.AccessTokenHeader, "Expected ACCESS_TOKEN_HEADER to be set to default 'Authorization'")
	assert.Equal(t, "https://example.com/.well-known/jwks.json", config.JwksUrl, "Expected JWKS_URL to be set to 'https://example.com/.well-known/jwks.json'")
}

func TestLoadConfig_MissingRequired(t *testing.T) {
	// Unset environment variables to simulate missing required ones
	os.Unsetenv("JWKS_URL")

	assert.Panics(t, func() {
		LoadConfig()
	}, "Expected LoadConfig to panic when JWKS_URL is missing")
}
