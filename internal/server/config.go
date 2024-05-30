// Package server manages server configuration and initialization based on environment variables.
// It includes utilities to load and handle configurations necessary for JWT authentication
// in a web server environment.
package server

import (
	"github.com/reMarkable/envconfig/v2"
	"log/slog"
)

// Config holds the configuration necessary for the server to validate JWT tokens.
// It extracts these settings from environment variables, providing defaults and
// enforcing required settings as needed.
type Config struct {
	AccessTokenHeader string `envconfig:"ACCESS_TOKEN_HEADER" default:"Authorization"` // AccessTokenHeader specifies the HTTP header name where the access token is expected. Default is "Authorization".
	JwksUrl           string `envconfig:"JWKS_URL" required:"true"`                    // JwksUrl is the URL pointing to the JSON Web Key Set (JWKS) resource for validating JWT tokens. This field is required.
}

// LoadConfig initializes a Config object from environment variables using the envconfig package.
// It logs a fatal error and exits the application if there is an error during loading,
// ensuring that the server does not start with invalid or incomplete configuration settings.
func LoadConfig() Config {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		slog.Error("Failed to load environment variables", "error", err)
		panic(err)
	}
	return config
}
