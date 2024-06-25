package handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

var correlationIdKey = contextKey("correlation_id")

func CorrelationId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrId := uuid.New().String()
		ctx := context.WithValue(r.Context(), correlationIdKey, corrId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Public getter for correlation ID
func GetCorrelationId(r *http.Request) (string, bool) {
	corrId, ok := r.Context().Value(correlationIdKey).(string)
	return corrId, ok
}
