package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
)

var (
	jwks      jwk.Set
	jwksMutex sync.RWMutex
)

// UpdateJWKS periodically fetches the JWKS from the specified URL and updates the cached copy.
// It runs as a goroutine, starting the update immediately upon launch and then at the specified interval.
func UpdateJWKS(url string, refreshInterval time.Duration) {
	fetchAndUpdateJWKS := func() {
		fetchedJWKS, err := jwk.Fetch(context.Background(), url)
		if err != nil {
			fmt.Printf("Error fetching JWKS: %v\n", err)
			return
		}
		jwksMutex.Lock()
		jwks = fetchedJWKS
		jwksMutex.Unlock()
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
