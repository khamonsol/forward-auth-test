package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strings"
)

const authRequestTxIdKey = "AUTH_REQUEST_TX_ID"

func AuthRequestTransactionIdHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := strings.Replace(uuid.New().String(), "-", "", -1) // Generate a new UUID for the correlation ID
		ctx := context.WithValue(r.Context(), authRequestTxIdKey, correlationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetAuthRequestTxId(r *http.Request) string {
	ctx := r.Context()
	val, ok := ctx.Value(authRequestTxIdKey).(string)
	if !ok {
		slog.Error("Unable to get correlation ID")
	}
	return val
}
