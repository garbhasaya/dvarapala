package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dvarapala/internal/platform/auth"
	"dvarapala/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestRouterAuthentication(t *testing.T) {
	jwtManager := auth.NewJWTManager("secret", 1*time.Hour)
	// We don't need a real service for this test as we just want to see if middleware blocks
	userHandler := user.NewHandler(nil) 
	router := NewRouter(userHandler, jwtManager)

	tests := []struct {
		name           string
		method         string
		url            string
		wantStatusCode int
	}{
		{"Health public", "GET", "/health", http.StatusOK},
		{"Users Auth public", "POST", "/users/auth", http.StatusBadRequest}, // 400 because of empty body
		{"Users List protected", "GET", "/users", http.StatusUnauthorized},
		{"Users Create protected", "POST", "/users", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body *bytes.Buffer
			if tt.method == "POST" {
				body = bytes.NewBuffer([]byte("{}"))
			}
			
			var req *http.Request
			if body != nil {
				req, _ = http.NewRequest(tt.method, tt.url, body)
			} else {
				req, _ = http.NewRequest(tt.method, tt.url, nil)
			}
			
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatusCode, rr.Code)
		})
	}
}

func TestRouterAuthentication_ValidToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("secret", 1*time.Hour)
	userHandler := user.NewHandler(nil) 
	router := NewRouter(userHandler, jwtManager)

	token, _ := jwtManager.Generate(1)

	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// If middleware works, it shouldn't be StatusUnauthorized (401).
	// It will be 500 because s.List is nil and it panics, and Recoverer catches it.
	assert.NotEqual(t, http.StatusUnauthorized, rr.Code)
}
