package middleware

import (
	"context"
	"fmt"
	"github.com/SoleaEnergy/forwardAuth/internal/util"
	"log/slog"
	"net/http"
	"strings"
)

const rolePolicyCheckKey = "ROLE_POLICY_CHECK"

// RoleVerificationHandler checks if the JWT contains any of the roles required to access
// the endpoint specified in the policy. This handler will continue even if the check fails
func RoleVerificationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var newCtx context.Context
		correlationId := GetAuthRequestTxId(r)
		token := ParseToken(w, r)
		if token == nil {
			newCtx = context.WithValue(r.Context(), rolePolicyCheckKey, false)
			next.ServeHTTP(w, r.WithContext(newCtx))
			return
		}

		claims := GetClaims(w, r, token)
		if claims == nil {
			newCtx = context.WithValue(r.Context(), rolePolicyCheckKey, false)
			next.ServeHTTP(w, r.WithContext(newCtx))
			return
		}

		roles := GetRoles(w, r, claims)
		if roles == nil {
			newCtx = context.WithValue(r.Context(), rolePolicyCheckKey, false)
			next.ServeHTTP(w, r.WithContext(newCtx))
			return
		}

		accessPolicy, err := GetPolicyFromRequest(r)
		if err != nil {
			msg := fmt.Sprintf("Error getting policy from request: %s", err.Error())
			util.HandleError(w, msg, http.StatusInternalServerError, correlationId)
			return
		}

		if !contains((*roles), accessPolicy.Roles) {
			msg := fmt.Sprintf("User has no roles in %s", strings.Join(accessPolicy.Roles, ","))
			slog.Info(msg)
			// Store the policy in the context of the request.
			newCtx = context.WithValue(r.Context(), rolePolicyCheckKey, false)
		} else {
			newCtx = context.WithValue(r.Context(), rolePolicyCheckKey, true)
		}
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
func GetRolePolicyStatusFromContext(r *http.Request) bool {
	ctx := r.Context()
	val, ok := ctx.Value(rolePolicyCheckKey).(bool)
	if !ok {
		slog.Error("User role check failed")
	}
	return val
}

// contains checks if any of the roles in the JWT are found in the required roles list.
func contains(jwtRoles []string, requiredRoles []string) bool {
	for _, role := range jwtRoles {
		strRole := role
		for _, reqRole := range requiredRoles {
			if reqRole == strRole {
				return true
			}
		}
	}
	return false
}
