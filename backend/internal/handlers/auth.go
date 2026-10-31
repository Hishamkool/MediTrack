package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"meditrack/internal/auth"
	"meditrack/internal/middleware"
	"meditrack/internal/models"
	"meditrack/internal/repository"
)

type AuthHandler struct {
	Users  *repository.UserRepo
	Tokens *auth.TokenService
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      *models.User `json:"user"`
}

// Login: POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		ValidationError(w, map[string]string{"email": "email and password are required"})
		return
	}

	user, err := h.Users.GetByEmail(r.Context(), email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		ServerError(w, r, err)
		return
	}
	// Same message for "no such user" and "wrong password", so attackers
	// can't find out which emails exist.
	if user == nil || !auth.CheckPassword(user.PasswordHash, req.Password) {
		Error(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if !user.IsActive {
		Error(w, http.StatusForbidden, "this account has been deactivated")
		return
	}

	token, expires, err := h.Tokens.Issue(user.ID)
	if err != nil {
		ServerError(w, r, err)
		return
	}
	JSON(w, http.StatusOK, loginResponse{Token: token, ExpiresAt: expires, User: user})
}

// Me: GET /api/auth/me — returns the logged-in user.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, middleware.UserFromContext(r.Context()))
}
