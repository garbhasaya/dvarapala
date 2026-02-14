package http

import (
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
		{"Health unauthenticated", "GET", "/health", http.StatusUnauthorized},
		{"Users unauthenticated", "GET", "/users", http.StatusUnauthorized},
		{"Users Auth unauthenticated", "POST", "/users/auth", http.StatusUnauthorized},
		{"Swagger unauthenticated", "GET", "/swagger/index.html", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, nil)
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

	req, _ := http.NewRequest("GET", "/health", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// It should NOT be Unauthorized. It might be 500 or 404 depending on the handler's dependencies,
	// but it shouldn't be 401. 
	// /health handler should return 200 OK.
	assert.Equal(t, http.StatusOK, rr.Code)
}
