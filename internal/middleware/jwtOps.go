package middleware

import (
	"context"
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/server"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
)

var (
	jwks      jwk.Set
	jwksMutex sync.RWMutex
)

func GetRoles(w http.ResponseWriter, r *http.Request, claims *jwt.MapClaims) *[]string {
	correlationId := GetAuthRequestTxId(r)
	var roles []interface{}
	if (*claims)["roles"] != nil {
		var ok bool
		roles, ok = (*claims)["roles"].([]interface{})
		if !ok {
			msg := fmt.Sprintf("Error getting roles from claims: %s", (*claims)["roles"])
			util.HandleError(w, msg, http.StatusInternalServerError, correlationId)
			return nil
		}
	}
	var roleList []string
	for _, role := range roles {
		roleList = append(roleList, fmt.Sprintf("%s", role))
	}
	return &roleList
}

func GetClaims(w http.ResponseWriter, r *http.Request, token *jwt.Token) *jwt.MapClaims {
	correlationId := GetAuthRequestTxId(r)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		msg := "Error getting claims from token"
		util.HandleError(w, msg, http.StatusInternalServerError, correlationId)
		return nil
	}
	return &claims
}

func ParseToken(w http.ResponseWriter, r *http.Request) *jwt.Token {
	sCfg := server.LoadConfig()
	tokenString := r.Header.Get(sCfg.AccessTokenHeader)
	correlationId := GetAuthRequestTxId(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Return a dummy key since the token has already been validated.
		return []byte("secret"), nil
	})

	// Should never happen because this token should have already been validated in an earlier step
	if err != nil || !token.Valid {
		util.HandleError(w, "Unauthorized: Invalid token", http.StatusInternalServerError, correlationId)
		return nil
	}
	return token
}

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
