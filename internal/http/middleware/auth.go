package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/pliffdax/sparrow-api/internal/security"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		raw := strings.TrimSpace(h[len("Bearer "):])
		t, claims, err := security.Parse(raw)
		if err != nil || !t.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}
		sub, _ := claims["sub"].(string)
		if sub == "" {
			http.Error(w, "invalid token payload", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
