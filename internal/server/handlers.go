package server

import (
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/policy"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	// Respond with OK status to indicate successful authentication
	w.WriteHeader(http.StatusOK)
}

// UserVerificationHandler checks if the JWT's preferred_username matches any of the configured
// users in the policy. If a match is found, it allows the request and skips further checks.
func UserVerificationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Omitted: return the public key for signature validation
			return nil, nil
		})

		userPolicy, err := policy.GetPolicyFromRequest(r)
		if err != nil {
			msg := fmt.Sprintf("Error getting policy from request: %s", err.Error())
			util.HandleError(w, msg, http.StatusInternalServerError)
		}

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username := claims["preferred_username"].(string)
		localPart := strings.Split(username, "@")[0]
		for _, user := range userPolicy.Users {
			if user == localPart {
				// User is verified, pass the request along with no further checks
				next.ServeHTTP(w, r)
				return
			}
		}

		// If no user match, continue to the next handler (role verification)
		next.ServeHTTP(w, r)
	})
}

// RoleVerificationHandler checks if the JWT contains any of the roles required to access
// the endpoint specified in the policy.
func RoleVerificationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Omitted: return the public key for signature validation
			return nil, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		roles := claims["roles"].([]interface{})
		accessPolicy, err := policy.GetPolicyFromRequest(r)
		if err != nil {
			util.HandleError(w, err.Error(), http.StatusInternalServerError)
		}

		if !contains(roles, accessPolicy.Roles) {
			msg := fmt.Sprintf("User has no roles in %s", strings.Join(accessPolicy.Roles, ","))
			util.HandleError(w, msg, http.StatusUnauthorized)
			return
		}

		// If roles match, pass the request along
		next.ServeHTTP(w, r)
	})
}

// contains checks if any of the roles in the JWT are found in the required roles list.
func contains(jwtRoles []interface{}, requiredRoles []string) bool {
	for _, role := range jwtRoles {
		strRole := role.(string)
		for _, reqRole := range requiredRoles {
			if reqRole == strRole {
				return true
			}
		}
	}
	return false
}

// JwtValidationHandler returns http.Handler that validates JWT tokens using the cached JWKS.
// It verifies that the token has not been altered by checking its signature against the JWKS keys.
func JwtValidationHandler(next http.Handler, config *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get(config.AccessTokenHeader)
		if tokenString == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		jwksMutex.RLock()
		localJWKS := jwks
		jwksMutex.RUnlock()

		if localJWKS == nil {
			util.HandleError(w, "Internal Server Error: JWKS not loaded", http.StatusInternalServerError)
			return
		}

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			keyID, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("expected jwt to have 'kid' header")
			}
			key, ok := localJWKS.LookupKeyID(keyID)
			if !ok {
				return nil, fmt.Errorf("unable to find key %q", keyID)
			}
			var pubKey interface{}
			if err := key.Raw(&pubKey); err != nil {
				return nil, fmt.Errorf("unable to get raw key for keyID %q: %w", keyID, err)
			}
			return pubKey, nil
		})

		if err != nil {
			msg := fmt.Sprintf("JWT validation error: %v", err)
			util.HandleError(w, msg, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
