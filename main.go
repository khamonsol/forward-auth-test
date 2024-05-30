// Package forward_auth provides functionality for handling forward authentication with Traefik.
package forward_auth

import (
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/policy"
	"github.com/SoleaEnergy/forwardAuth/internal/server"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// main is the entry point of the forward auth server.
func main() {
	// Load server configuration
	config := server.LoadConfig()

	// Start JWKS update routine
	go server.UpdateJWKS(config.JwksUrl, 24*time.Hour) // Adjust refresh interval as needed

	// Define the main handler that simply returns HTTP 200 OK
	mainHandler := http.HandlerFunc(server.SuccessHandler)

	// Set up the middleware chain
	// Start with the innermost handler (mainHandler) and wrap it with each middleware going outward
	handlerChain := server.RoleVerificationHandler(mainHandler)       // Wrap main handler with role verification
	handlerChain = server.UserVerificationHandler(handlerChain)       // Wrap with user verification
	handlerChain = server.JwtValidationHandler(handlerChain, &config) // Wrap with JWT validation
	handlerChain = policy.PolicyLoader(handlerChain)                  // Wrap with policy loader as the outermost middleware

	//Execution order: PolicyLoader -> JwtValidationHandler -> UserVerificationHandler -> RoleVerificationHandler -> SuccessHandler

	// Setup HTTP server and routes
	http.Handle("/", handlerChain) // You can specify more routes as needed

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		msg := fmt.Sprintf("Error starting server: %s", err)
		slog.Error(msg, err)
		os.Exit(1)
	}
}
