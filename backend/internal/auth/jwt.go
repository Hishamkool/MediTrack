package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService creates and checks signed access tokens.
type TokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService(secret string, ttl time.Duration) *TokenService {
	return &TokenService{secret: []byte(secret), ttl: ttl}
}

// Issue creates a token for the given user ID. The token only carries the
// ID; the role is looked up fresh from the database on every request, so a
// role change or deactivation takes effect immediately.
func (s *TokenService) Issue(userID int64) (string, time.Time, error) {
	expires := time.Now().Add(s.ttl)
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expires),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	return signed, expires, err
}

// Parse validates the token (signature + expiry) and returns the user ID.
func (s *TokenService) Parse(tokenString string) (int64, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, errors.New("invalid token subject")
	}
	return id, nil
}
