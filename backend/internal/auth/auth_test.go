package auth

import (
	"testing"
	"time"
)

func TestPasswordHashAndCheck(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(hash, "secret123") {
		t.Error("correct password was rejected")
	}
	if CheckPassword(hash, "wrong") {
		t.Error("wrong password was accepted")
	}
}

func TestTokenRoundTrip(t *testing.T) {
	s := NewTokenService("test-secret-that-is-long-enough-123456", time.Hour)
	token, _, err := s.Issue(42)
	if err != nil {
		t.Fatal(err)
	}
	id, err := s.Parse(token)
	if err != nil || id != 42 {
		t.Fatalf("got id=%d err=%v, want 42", id, err)
	}
}

func TestTokenRejectedWithWrongSecret(t *testing.T) {
	token, _, _ := NewTokenService("secret-one-secret-one-secret-one-xx", time.Hour).Issue(1)
	if _, err := NewTokenService("secret-two-secret-two-secret-two-xx", time.Hour).Parse(token); err == nil {
		t.Error("token signed with another secret was accepted")
	}
}

func TestExpiredTokenRejected(t *testing.T) {
	s := NewTokenService("test-secret-that-is-long-enough-123456", -time.Minute)
	token, _, _ := s.Issue(1)
	if _, err := s.Parse(token); err == nil {
		t.Error("expired token was accepted")
	}
}
