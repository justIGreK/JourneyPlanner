package handler

import (
	"context"
	"net/http"
	"strings"
)

type ContextKey string

const (
	UserLoginKey ContextKey = "user_id"
)

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		payload, err := h.User.ValidatePasetoToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserLoginKey, payload.UserLogin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
