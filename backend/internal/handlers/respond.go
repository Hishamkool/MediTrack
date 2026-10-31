// Package handlers contains the HTTP handlers (one file per resource).
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

// JSON writes v as a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// Error writes the standard error shape: {"error": "message"}.
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

// ValidationError writes 422 with per-field messages, e.g.
// {"error": "validation failed", "fields": {"email": "is required"}}.
func ValidationError(w http.ResponseWriter, fields map[string]string) {
	JSON(w, http.StatusUnprocessableEntity, map[string]any{
		"error":  "validation failed",
		"fields": fields,
	})
}

// ServerError logs the real error and sends a generic message, so internal
// details never leak to the client.
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal error", "method", r.Method, "path", r.URL.Path, "err", err)
	Error(w, http.StatusInternalServerError, "something went wrong")
}

// decodeJSON reads the request body into dst. Unknown fields are rejected
// so typos in the frontend are noticed early. Bodies are capped at 1 MB.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			Error(w, http.StatusBadRequest, "request body is empty")
		} else {
			Error(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		}
		return false
	}
	return true
}
