package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"meditrack/internal/models"
)

// Fakes so we can test without a database.
type fakeTokens map[string]int64

func (f fakeTokens) Parse(token string) (int64, error) {
	if id, ok := f[token]; ok {
		return id, nil
	}
	return 0, errors.New("bad token")
}

type fakeUsers map[int64]*models.User

func (f fakeUsers) GetByID(_ context.Context, id int64) (*models.User, error) {
	if u, ok := f[id]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

var ok = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func TestAuthenticate(t *testing.T) {
	tokens := fakeTokens{"good": 1, "inactive": 2}
	users := fakeUsers{
		1: {ID: 1, Role: models.RoleStaff, IsActive: true},
		2: {ID: 2, Role: models.RoleStaff, IsActive: false},
	}
	h := Authenticate(tokens, users)(ok)

	tests := []struct {
		name   string
		header string
		want   int
	}{
		{"no header", "", http.StatusUnauthorized},
		{"not bearer", "Basic abc", http.StatusUnauthorized},
		{"bad token", "Bearer nope", http.StatusUnauthorized},
		{"inactive user", "Bearer inactive", http.StatusUnauthorized},
		{"valid", "Bearer good", http.StatusOK},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Errorf("got %d, want %d", rec.Code, tc.want)
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	h := RequireRole(models.RoleAdmin)(ok)

	tests := []struct {
		name string
		user *models.User
		want int
	}{
		{"no user", nil, http.StatusUnauthorized},
		{"staff", &models.User{Role: models.RoleStaff}, http.StatusForbidden},
		{"admin", &models.User{Role: models.RoleAdmin}, http.StatusOK},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tc.user != nil {
				req = req.WithContext(WithUser(req.Context(), tc.user))
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Errorf("got %d, want %d", rec.Code, tc.want)
			}
		})
	}
}
