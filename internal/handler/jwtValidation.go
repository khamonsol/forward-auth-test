package handler

import (
	"context"
	"github.com/SoleaEnergy/forwardAuth/internal/policy"
	"github.com/SoleaEnergy/forwardAuth/internal/provider"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

// Private keys for context
var (
	isAllowedKey    = contextKey("isAllowed")
	errorStatusKey  = contextKey("errorStatus")
	errorMessageKey = contextKey("errorMessage")
	jwkCache        = jwk.NewCache(context.Background())
)

func newJWKSet(jwkUrl string) (jwk.Set, error) {
	jwkCache := jwk.NewCache(context.Background())

	// register a minimum refresh interval for this URL.
	// when not specified, defaults to Cache-Control and similar resp headers
	err := jwkCache.Register(jwkUrl, jwk.WithMinRefreshInterval(10*time.Minute))
	if err != nil {
		panic("failed to register jwk location")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// fetch once on application startup
	_, err = jwkCache.Refresh(ctx, jwkUrl)
	if err != nil {
		return nil, err
	}
	// create the cached key set
	return jwk.NewCachedSet(jwkCache, jwkUrl), nil
}

func JwtValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract the token from the request headers
		tokenString := extractToken(r)
		if tokenString == "" {
			r = setErrorStatus(r, http.StatusUnauthorized)
			r = setErrorMessage(r, "Unauthorized: Missing token")
			next.ServeHTTP(w, r)
			return
		}

		// Create Kubernetes client
		kubeClient, err := util.NewKubernetesClient()
		if err != nil {
			r = setErrorStatus(r, http.StatusInternalServerError)
			r = setErrorMessage(r, "Internal Server Error: Failed to create Kubernetes client")
			next.ServeHTTP(w, r)
			return
		}

		// Retrieve the policy for the requested resource
		policy, err := policy.LoadPolicies(r.Host, kubeClient)
		if err != nil {
			r = setErrorStatus(r, http.StatusInternalServerError)
			r = setErrorMessage(r, "Internal Server Error: Failed to load policies")
			next.ServeHTTP(w, r)
			return
		}
		err = policy.GetPolicy(r.URL.Path, r.Method)
		if err != nil {
			r = setErrorStatus(r, http.StatusForbidden)
			r = setErrorMessage(r, "Forbidden: Policy not found")
			next.ServeHTTP(w, r)
			return
		}

		// Load provider configuration
		providerConfig := provider.NewConfig(policy.ProviderType, policy.ProviderName, "default")
		err = providerConfig.LoadProviderConfig(kubeClient)
		if err != nil {
			r = setErrorStatus(r, http.StatusInternalServerError)
			r = setErrorMessage(r, "Internal Server Error: Failed to load provider config")
			next.ServeHTTP(w, r)
			return
		}

		// Validate the token
		if !validateToken(tokenString, providerConfig, policy) {
			r = setErrorStatus(r, http.StatusUnauthorized)
			r = setErrorMessage(r, "Unauthorized: Invalid token")
			next.ServeHTTP(w, r)
			return
		}

		// Mark the request as allowed
		r = setAllowed(r, true)
		next.ServeHTTP(w, r)
	})
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

func validateToken(tokenString string, config provider.Provider, policy *policy.Policy) bool {
	jwksURL := config.GetIssuerUrl() + "/.well-known/jwks.json"
	keySet, err := newJWKSet(jwksURL)
	if err != nil {
		return false
	}

	token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(keySet))
	if err != nil {
		return false
	}
	// Extract roles from the token claims
	roles, ok := token.Get("roles")
	if !ok {
		return false
	}

	rolesSlice, ok := roles.([]interface{})
	if !ok {
		return false
	}

	// Check roles in the token against the policy
	for _, role := range rolesSlice {
		for _, policyRole := range policy.Roles {
			if role == policyRole {
				return true
			}
		}
	}

	return false
}

// Private setters
func setAllowed(r *http.Request, allowed bool) *http.Request {
	ctx := context.WithValue(r.Context(), isAllowedKey, allowed)
	return r.WithContext(ctx)
}

func setErrorStatus(r *http.Request, status int) *http.Request {
	ctx := context.WithValue(r.Context(), errorStatusKey, status)
	return r.WithContext(ctx)
}

func setErrorMessage(r *http.Request, message string) *http.Request {
	ctx := context.WithValue(r.Context(), errorMessageKey, message)
	return r.WithContext(ctx)
}

// Public getters
func GetAllowed(r *http.Request) (bool, bool) {
	allowed, ok := r.Context().Value(isAllowedKey).(bool)
	return allowed, ok
}

func GetErrorStatus(r *http.Request) (int, bool) {
	status, ok := r.Context().Value(errorStatusKey).(int)
	return status, ok
}

func GetErrorMessage(r *http.Request) (string, bool) {
	message, ok := r.Context().Value(errorMessageKey).(string)
	return message, ok
}
