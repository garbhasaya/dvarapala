package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"dvarapala/internal/platform/render"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

// Middleware returns a middleware that authenticates requests using JWT.
func Middleware(manager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				slog.Warn("missing authorization header", "path", r.URL.Path, "remote_addr", r.RemoteAddr)
				render.Error(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				slog.Warn("invalid authorization header format", "path", r.URL.Path, "remote_addr", r.RemoteAddr)
				render.Error(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			token := parts[1]
			claims, err := manager.Verify(token)
			if err != nil {
				slog.Warn("invalid or expired token", "path", r.URL.Path, "remote_addr", r.RemoteAddr, "error", err)
				render.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
