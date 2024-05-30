package policy

import (
	"context"
	"errors"
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"net/http"
)

// policyHeaderKey is the context key used for storing the policy in the request context.
// It is used to retrieve the policy data within downstream handlers.
const policyHeaderKey = "POLICY"

// PolicyLoader is an HTTP middleware that intercepts the request to load and attach
// the appropriate policy configuration based on the request details. It performs
// three main actions: it initializes a new configuration, loads the specific policy
// based on the request's host, and then fetches the specific policy for the requested
// path and method. If any step fails, it sends an HTTP error response.

func PolicyLoader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Initialize and load the configuration based on the host of the request.
		policyConfig, err := NewConfig()
		if err != nil {
			msg := fmt.Sprintf("Failed to load config map error: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		err = policyConfig.LoadConfig(r.Host)
		if err != nil {
			msg := fmt.Sprintf("Failed to load policy configuration for host %s error: %v", r.Host, err)
			util.HandleError(w, msg, http.StatusInternalServerError)
			return
		}

		// Retrieve the policy using the request's path and method.
		policy, err := policyConfig.GetPolicy(r.URL.Path, r.Method)
		if err != nil {
			msg := fmt.Sprintf("Failed to load policy for path %s method: %v", r.URL.Path, r.Method)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		// Store the policy in the context of the request.
		ctx := context.WithValue(r.Context(), policyHeaderKey, policy)

		// Proceed with the next handler, passing the context with the loaded policy.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPolicyFromRequest extracts the policy data from the request's context.
// This function is intended to be used by downstream handlers that need to access
// the policy associated with the request. It assumes that the policy data exists
// in the context; otherwise, it will return an error.
func GetPolicyFromRequest(r *http.Request) (Policy, error) {
	ctx := r.Context()
	val, ok := ctx.Value(policyHeaderKey).(*Policy)
	if !ok {
		return Policy{}, errors.New("policy data not found in request context; ensure PolicyLoader middleware is configured properly")
	}
	return *val, nil
}
