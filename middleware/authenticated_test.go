package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name          string
		setupRequest  func() *http.Request
		expectedToken string
		expectError   bool
	}{
		{
			name: "Valid token in header",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/user/list", nil)
				req.Header.Set("Authorization", "Bearer valid-token")
				return req
			},
			expectedToken: "valid-token",
			expectError:   false,
		},
		{
			name: "Valid token in query param",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/user/list", nil)
				q := req.URL.Query()
				q.Add("authorization", "Bearer valid-token")
				req.URL.RawQuery = q.Encode()
				return req
			},
			expectedToken: "valid-token",
			expectError:   false,
		},
		{
			name: "Missing token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/user/list", nil)
				return req
			},
			expectedToken: "",
			expectError:   true,
		},
		{
			name: "Invalid header format",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/user/list", nil)
				req.Header.Set("Authorization", "invalid-format")
				return req
			},
			expectedToken: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupRequest()
			token, err := extractToken(req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET_KEY")
	defer func() {
		os.Setenv("JWT_SECRET_KEY", originalSecret)
	}()

	// Set a test secret key
	testSecret := "test-secret-key"
	os.Setenv("JWT_SECRET_KEY", testSecret)

	// Create a valid token
	validClaims := jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)
	validTokenString, err := validToken.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("Error creating test token: %v", err)
	}

	// Create an expired token
	expiredClaims := jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, err := expiredToken.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("Error creating expired token: %v", err)
	}

	tests := []struct {
		name        string
		setupEnv    func()
		tokenString string
		expectError bool
	}{
		{
			name:        "Valid token",
			setupEnv:    func() {},
			tokenString: validTokenString,
			expectError: false,
		},
		{
			name:        "Expired token",
			setupEnv:    func() {},
			tokenString: expiredTokenString,
			expectError: true,
		},
		{
			name: "Missing secret key",
			setupEnv: func() {
				os.Setenv("JWT_SECRET_KEY", "")
			},
			tokenString: validTokenString,
			expectError: true,
		},
		{
			name:        "Invalid token",
			setupEnv:    func() {},
			tokenString: "invalid-token",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			claims, err := validateToken(tt.tokenString)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, "testuser", claims["username"])
			}

			// Reset environment
			os.Setenv("JWT_SECRET_KEY", testSecret)
		})
	}
}
