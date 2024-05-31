package middleware

import (
	"net/http"
)

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	// Respond with OK status to indicate successful authentication
	hasValidUser := GetUserPolicyStatusFromContext(r)
	hasValidRole := GetRolePolicyStatusFromContext(r)
	if hasValidUser || hasValidRole {
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusUnauthorized)
}
