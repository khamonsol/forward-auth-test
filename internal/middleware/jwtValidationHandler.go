package middleware

import (
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/server"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

// JwtValidationHandler returns http.Handler that validates JWT tokens using the cached JWKS.
// It verifies that the token has not been altered by checking its signature against the JWKS keys.
// This handler will stop further processing because there is no reason to continue if the token is
// not trusted.
func JwtValidationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sCfg := server.LoadConfig()
		tokenString := r.Header.Get(sCfg.AccessTokenHeader)
		correlationID := GetAuthRequestTxId(r)
		if tokenString == "" {
			util.HandleError(w, "Unauthorized: No token provided", http.StatusUnauthorized, correlationID)
			return
		}

		jwksMutex.RLock()
		localJWKS := jwks
		jwksMutex.RUnlock()

		if localJWKS == nil {
			util.HandleError(w, "Internal Server Error: JWKS not loaded", http.StatusInternalServerError, correlationID)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
		}, jwt.WithValidMethods(strings.Split(sCfg.ValidAlgs, ",")))

		if err != nil {
			msg := fmt.Sprintf("JWT validation error: %v", err)
			util.HandleError(w, msg, http.StatusUnauthorized, correlationID)
			return
		}

		if !token.Valid {
			util.HandleError(w, "Unauthorized: Invalid token", http.StatusUnauthorized, correlationID)
			return
		}

		// Token is valid, pass the request along
		next.ServeHTTP(w, r)
	})
}
