package handlers

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"meditrack/internal/auth"
	"meditrack/internal/middleware"
	"meditrack/internal/models"
	"meditrack/internal/repository"
)

const minPasswordLen = 8

// UserHandler serves the admin-only user management endpoints.
type UserHandler struct {
	Users *repository.UserRepo
}

// List: GET /api/users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.Users.List(r.Context())
	if err != nil {
		ServerError(w, r, err)
		return
	}
	JSON(w, http.StatusOK, users)
}

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // optional, defaults to "staff"
}

// Create: POST /api/users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Role == "" {
		req.Role = models.RoleStaff
	}

	fields := map[string]string{}
	if req.Name == "" {
		fields["name"] = "is required"
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		fields["email"] = "must be a valid email"
	}
	if len(req.Password) < minPasswordLen {
		fields["password"] = "must be at least 8 characters"
	}
	if !validRole(req.Role) {
		fields["role"] = "must be 'admin' or 'staff'"
	}
	if len(fields) > 0 {
		ValidationError(w, fields)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		ServerError(w, r, err)
		return
	}
	user := &models.User{
		Name: req.Name, Email: req.Email, PasswordHash: hash,
		Role: req.Role, IsActive: true,
	}
	if err := h.Users.Create(r.Context(), user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			ValidationError(w, map[string]string{"email": "is already in use"})
			return
		}
		ServerError(w, r, err)
		return
	}
	JSON(w, http.StatusCreated, user)
}

// updateUserRequest uses pointers so we can tell "not sent" (nil)
// apart from "sent as empty/false".
type updateUserRequest struct {
	Name     *string `json:"name"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
	Password *string `json:"password"` // admin resets a user's password
}

// Update: PATCH /api/users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid user id")
		return
	}
	var req updateUserRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	user, err := h.Users.GetByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		Error(w, http.StatusNotFound, "user not found")
		return
	} else if err != nil {
		ServerError(w, r, err)
		return
	}

	// An admin may not deactivate or demote themselves. Since only admins
	// reach this handler, this guarantees at least one active admin remains.
	me := middleware.UserFromContext(r.Context())
	if me.ID == user.ID {
		if (req.IsActive != nil && !*req.IsActive) || (req.Role != nil && *req.Role != models.RoleAdmin) {
			Error(w, http.StatusBadRequest, "you cannot deactivate or demote your own account")
			return
		}
	}

	fields := map[string]string{}
	if req.Name != nil {
		if name := strings.TrimSpace(*req.Name); name == "" {
			fields["name"] = "cannot be empty"
		} else {
			user.Name = name
		}
	}
	if req.Role != nil {
		if !validRole(*req.Role) {
			fields["role"] = "must be 'admin' or 'staff'"
		} else {
			user.Role = *req.Role
		}
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.Password != nil {
		if len(*req.Password) < minPasswordLen {
			fields["password"] = "must be at least 8 characters"
		} else if user.PasswordHash, err = auth.HashPassword(*req.Password); err != nil {
			ServerError(w, r, err)
			return
		}
	}
	if len(fields) > 0 {
		ValidationError(w, fields)
		return
	}

	if err := h.Users.Update(r.Context(), user); err != nil {
		ServerError(w, r, err)
		return
	}
	JSON(w, http.StatusOK, user)
}

func validRole(role string) bool {
	return role == models.RoleAdmin || role == models.RoleStaff
}
