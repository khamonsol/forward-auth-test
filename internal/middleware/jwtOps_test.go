package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/assert"
)

var testJWKS = `{"keys":[{"kty":"RSA","n":"someModulus","e":"AQAB","kid":"test-key-id"}]}`

// MockJWKSHandler is an HTTP handler that returns a mocked JWKS.
func MockJWKSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(testJWKS))
}

func TestUpdateJWKS(t *testing.T) {
	// Start a local HTTP server to serve the mock JWKS
	server := httptest.NewServer(http.HandlerFunc(MockJWKSHandler))
	defer server.Close()

	// Override fetchAndUpdateJWKS to notify the wait group when the initial fetch is done
	fetchAndUpdateJWKS := func(url string) {
		fetchedJWKS, err := jwk.Fetch(context.Background(), url)
		if err != nil {
			fmt.Printf("Error fetching JWKS: %v\n", err)
			return
		}
		jwksMutex.Lock()
		jwks = fetchedJWKS
		jwksMutex.Unlock()
	}

	// Set up a wait group to wait for the initial fetch
	var wg sync.WaitGroup
	wg.Add(1)

	// Wrap the UpdateJWKS function to use the overridden fetchAndUpdateJWKS
	updateJWKS := func(url string, refreshInterval time.Duration) {
		fetchAndUpdateJWKS(url)
		wg.Done()

		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		for range ticker.C {
			fetchAndUpdateJWKS(url)
		}
	}

	// Start the UpdateJWKS function
	go updateJWKS(server.URL, time.Second*10)

	// Wait for the initial fetch to complete
	wg.Wait()

	// Verify that the jwks variable was updated
	jwksMutex.RLock()
	defer jwksMutex.RUnlock()

	assert.NotNil(t, jwks, "Expected jwks to be non-nil")
	assert.Equal(t, 1, jwks.Len(), "Expected jwks to contain 1 key")
	key, ok := jwks.LookupKeyID("test-key-id")
	assert.True(t, ok, "Expected jwks to contain a key with ID 'test-key-id'")
	assert.NotNil(t, key, "Expected the key to be non-nil")
}
