// Package middleware has the HTTP middleware that protects routes.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"meditrack/internal/models"
)

// TokenParser and UserGetter are small interfaces so tests can pass fakes
// instead of a real database.
type TokenParser interface {
	Parse(token string) (int64, error)
}

type UserGetter interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

// contextKey is a private type so no other package can clash with our key.
type contextKey struct{}

// UserFromContext returns the logged-in user set by Authenticate.
func UserFromContext(ctx context.Context) *models.User {
	u, _ := ctx.Value(contextKey{}).(*models.User)
	return u
}

// WithUser stores a user in the context (used by Authenticate and tests).
func WithUser(ctx context.Context, u *models.User) context.Context {
	return context.WithValue(ctx, contextKey{}, u)
}

// Authenticate requires a valid "Authorization: Bearer <token>" header.
// It loads the user from the database so deactivated users are rejected
// at once, even if their token has not expired yet.
func Authenticate(tokens TokenParser, users UserGetter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
				return
			}

			userID, err := tokens.Parse(token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			user, err := users.GetByID(r.Context(), userID)
			if err != nil || !user.IsActive {
				writeError(w, http.StatusUnauthorized, "account not found or inactive")
				return
			}

			next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
		})
	}
}

// RequireRole only lets users with one of the given roles through.
// It must be used after Authenticate.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())
			if user == nil {
				writeError(w, http.StatusUnauthorized, "not logged in")
				return
			}
			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			writeError(w, http.StatusForbidden, "you do not have permission to do this")
		})
	}
}

// writeError mirrors handlers.Error; duplicated here to avoid an import cycle.
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
