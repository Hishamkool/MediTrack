// Package auth handles password hashing and JWT access tokens.
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a bcrypt hash that is safe to store in the database.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword reports whether plain matches the stored hash.
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
