package forwardAuth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SoleaEnergy/forwardAuth/internal/middleware"
	"github.com/SoleaEnergy/forwardAuth/internal/server"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/assert"
)

// MockUpdateJWKS mocks the JWKS update function to avoid network calls.
func MockUpdateJWKS(url string, refreshInterval time.Duration) {
	// Mock implementation of JWKS fetching and updating
	fetchAndUpdateJWKS := func() {
		fetchedJWKS := jwk.NewSet()
		// Mock updating the global JWKS
		middleware.SetJWKS(fetchedJWKS)
	}

	// Initial fetch
	fetchAndUpdateJWKS()

	// Periodic updates
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		fetchAndUpdateJWKS()
	}
}

// MockLoadConfig returns a mock configuration for testing.
func MockLoadConfig() server.Config {
	return server.Config{
		JwksUrl: "https://example.com/.well-known/jwks.json",
	}
}

func TestSetupHandlerChain(t *testing.T) {
	// Mock the configuration
	config := MockLoadConfig()

	// Setup the middleware chain with the mock JWKS updater
	handler := setupHandlerChain(config, MockUpdateJWKS)

	// Create a test server
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Make a request to the test server
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestStartServer(t *testing.T) {

	// Run the server in a separate goroutine
	go func() {
		err := startServer()
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Allow some time for the server to start
	time.Sleep(1 * time.Second)

	// Make a request to the server
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
