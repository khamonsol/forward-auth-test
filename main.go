// Package forward_auth provides functionality for handling forward authentication with Traefik.
package forwardAuth

import (
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/middleware"
	"log/slog"
	"net/http"
	"os"
)

func setupHandlerChain() http.Handler {
	// Define the main handler that simply returns HTTP 200 OK
	mainHandler := middleware.SuccessHandler()

	// Set up the middleware chain
	// Start with the innermost handler (mainHandler) and wrap it with each middleware going outward
	handlerChain := middleware.JwtValidationHandler(mainHandler)            // Verify and authenticate
	handlerChain = middleware.SetupContextHandler(handlerChain)             // Setup context for validation
	handlerChain = middleware.AuthRequestTransactionIdHandler(handlerChain) // Add a transaction ID to the request
	return handlerChain
}

func startServer() error {
	http.Handle("/", setupHandlerChain())
	return http.ListenAndServe(":8080", nil)
}

// main is the entry point of the forward auth server.
func main() {
	err := startServer()
	if err != nil {
		msg := fmt.Sprintf("Error starting server: %s", err)
		slog.Error(msg, err)
		os.Exit(1)
	}
}
