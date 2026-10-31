// Package models holds the plain data structs shared across the app.
package models

import "time"

// Roles a user can have.
const (
	RoleAdmin = "admin"
	RoleStaff = "staff"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" means: never send this in JSON
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}
