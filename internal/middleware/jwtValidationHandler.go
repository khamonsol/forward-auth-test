package middleware

import (
	"github.com/SoleaEnergy/forwardAuth/internal/jwtops"
	"github.com/SoleaEnergy/forwardAuth/internal/provider"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"net/http"
)

// JwtValidationHandler returns http.Handler that validates JWT tokens using the cached JWKS.
// It verifies that the token has not been altered by checking its signature against the JWKS keys.
// This handler will stop further processing because there is no reason to continue if the token is
// not trusted.
func JwtValidationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrId := GetAuthRequestTxId(r)
		kApi := GetKubeApi(r)
		pol := GetPolicy(r)
		providerConfig := provider.NewConfig(pol.ProviderType, pol.ProviderName)
		if providerConfig != nil {
			util.HandleError(w, "Unable to construct provider config", http.StatusInternalServerError, corrId)
			return
		}
		err := providerConfig.LoadProviderConfig(kApi)
		if err != nil {
			util.HandleError(w, "Unable to load provider config", http.StatusInternalServerError, corrId)
			return
		}
		at, err := jwtops.LoadToken(r, providerConfig)
		if err != nil {
			util.HandleError(w, "Unable to load token", http.StatusInternalServerError, corrId)
			return
		}
		err = at.ValidateTokenForPolicy(*pol)
		if err != nil {
			util.HandleError(w, err.Error(), http.StatusUnauthorized, corrId)
			return
		}
		// Token is valid, pass the request along
		next.ServeHTTP(w, r)
	})
}
