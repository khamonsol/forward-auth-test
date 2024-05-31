package middleware

import (
	"context"
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"log/slog"
	"net/http"
	"strings"
)

const userPolicyCheckKey = "USER_POLICY_CHECK"

// UserVerificationHandler checks if the JWT's preferred_username matches any of the configured
// users in the policy. If a match is found, it allows the request and skips further checks.
func UserVerificationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var nxtCtx context.Context
		correlationID := GetAuthRequestTxId(r)

		token := ParseToken(w, r)
		if token == nil {
			nxtCtx = context.WithValue(r.Context(), userPolicyCheckKey, false)
			next.ServeHTTP(w, r.WithContext(nxtCtx))
			return
		}

		userPolicy, err := GetPolicyFromRequest(r)
		if err != nil {
			msg := fmt.Sprintf("Error getting policy from request: %s", err.Error())
			util.HandleError(w, msg, http.StatusInternalServerError, correlationID)
		}

		claims := GetClaims(w, r, token)
		if claims == nil {
			nxtCtx = context.WithValue(r.Context(), userPolicyCheckKey, false)
			next.ServeHTTP(w, r.WithContext(nxtCtx))
			return
		}

		username := (*claims)["preferred_username"].(string)
		localPart := strings.Split(username, "@")[0]
		for _, user := range userPolicy.Users {
			if user == localPart {
				ctx := context.WithValue(r.Context(), userPolicyCheckKey, true)
				// Proceed with the next handler, passing the context with the loaded policy.
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Store the policy in the context of the request.
		ctx := context.WithValue(r.Context(), userPolicyCheckKey, false)
		// Proceed with the next handler, passing the context with the loaded policy.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserPolicyStatusFromContext(r *http.Request) bool {
	ctx := r.Context()
	val, ok := ctx.Value(userPolicyCheckKey).(bool)
	if !ok {
		slog.Error("User role check failed")
	}
	return val
}
