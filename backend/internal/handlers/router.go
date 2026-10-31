package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"meditrack/internal/auth"
	"meditrack/internal/config"
	mw "meditrack/internal/middleware"
	"meditrack/internal/models"
	"meditrack/internal/repository"
)

// NewRouter wires every URL to its handler and middleware.
func NewRouter(cfg *config.Config, db *pgxpool.Pool) http.Handler {
	users := repository.NewUserRepo(db)
	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.TokenTTL)

	authH := &AuthHandler{Users: users, Tokens: tokens}
	userH := &UserHandler{Users: users}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger) // logs method, path, status — never bodies or headers
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	r.Get("/healthz", Healthz(db))

	r.Route("/api", func(r chi.Router) {
		// Public
		r.Post("/auth/login", authH.Login)

		// Everything below needs a valid token.
		r.Group(func(r chi.Router) {
			r.Use(mw.Authenticate(tokens, users))

			r.Get("/auth/me", authH.Me)

			// Admin only
			r.Group(func(r chi.Router) {
				r.Use(mw.RequireRole(models.RoleAdmin))
				r.Get("/users", userH.List)
				r.Post("/users", userH.Create)
				r.Patch("/users/{id}", userH.Update)
			})
		})
	})

	return r
}
