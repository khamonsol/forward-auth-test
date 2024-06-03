package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	err := os.Setenv("ACCESS_TOKEN_HEADER", "X-Access-Token")
	if err != nil {
		t.FailNow()
	}
	err = os.Setenv("JWKS_URL", "https://example.com/.well-known/jwks.json")
	if err != nil {
		t.FailNow()
	}

	// Ensure environment variables are unset after the test
	defer func() {
		err := os.Unsetenv("ACCESS_TOKEN_HEADER")
		if err != nil {
			t.FailNow()
		}
		err = os.Unsetenv("JWKS_URL")
		if err != nil {
			t.FailNow()
		}
	}()

	config := LoadConfig()

	assert.Equal(t, "X-Access-Token", config.AccessTokenHeader, "Expected ACCESS_TOKEN_HEADER to be set to 'X-Access-Token'")
	assert.Equal(t, "https://example.com/.well-known/jwks.json", config.JwksUrl, "Expected JWKS_URL to be set to 'https://example.com/.well-known/jwks.json'")
}

func TestLoadConfig_DefaultHeader(t *testing.T) {
	// Set only the required environment variable
	err := os.Setenv("JWKS_URL", "https://example.com/.well-known/jwks.json")
	if err != nil {
		t.FailNow()
	}
	// Make sure we get the default value for the token header
	err = os.Unsetenv("ACCESS_TOKEN_HEADER")
	if err != nil {
		t.FailNow()
	}

	// Ensure environment variables are unset after the test
	defer func() {
		err := os.Unsetenv("JWKS_URL")
		if err != nil {
			t.FailNow()
		}
	}()

	config := LoadConfig()

	assert.Equal(t, "Authorization", config.AccessTokenHeader, "Expected ACCESS_TOKEN_HEADER to be set to default 'Authorization'")
	assert.Equal(t, "https://example.com/.well-known/jwks.json", config.JwksUrl, "Expected JWKS_URL to be set to 'https://example.com/.well-known/jwks.json'")
}

func TestLoadConfig_MissingRequired(t *testing.T) {
	// Unset environment variables to simulate missing required ones
	err := os.Unsetenv("JWKS_URL")
	if err != nil {
		t.FailNow()
	}

	assert.Panics(t, func() {
		LoadConfig()
	}, "Expected LoadConfig to panic when JWKS_URL is missing")
}

func TestLoadConfig_Singleton(t *testing.T) {
	// Set environment variables for testing
	err := os.Setenv("ACCESS_TOKEN_HEADER", "X-Access-Token")
	if err != nil {
		t.FailNow()
	}
	err = os.Setenv("JWKS_URL", "https://example.com/.well-known/jwks.json")
	if err != nil {
		t.FailNow()
	}

	// Ensure environment variables are unset after the test
	defer func() {
		err := os.Unsetenv("ACCESS_TOKEN_HEADER")
		if err != nil {
			t.FailNow()
		}
		err = os.Unsetenv("JWKS_URL")
		if err != nil {
			t.FailNow()
		}
	}()

	config1 := LoadConfig()
	config2 := LoadConfig()

	assert.Equal(t, &config1, &config2, "Expected LoadConfig to return the same instance")
}
