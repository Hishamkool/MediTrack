package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Healthz is used by Render/Fly to check the API is alive and can reach the DB.
func Healthz(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "db unavailable"})
			return
		}
		JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
