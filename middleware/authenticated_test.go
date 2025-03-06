package middleware

import (
	"go-authentication-exercise/user/entity"
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

func TestGetUsernameFromJwt(t *testing.T) {
	tests := []struct {
		name           string
		claims         map[string]interface{}
		expectedResult string
		expectError    bool
	}{
		{
			name: "Valid username",
			claims: map[string]interface{}{
				"username": "testuser",
			},
			expectedResult: "testuser",
			expectError:    false,
		},
		{
			name: "Missing username",
			claims: map[string]interface{}{
				"other": "value",
			},
			expectedResult: "",
			expectError:    true,
		},
		{
			name: "Empty username",
			claims: map[string]interface{}{
				"username": "",
			},
			expectedResult: "",
			expectError:    true,
		},
		{
			name: "Non-string username",
			claims: map[string]interface{}{
				"username": 123,
			},
			expectedResult: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, err := getUsernameFromJwt(tt.claims)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, username)
			}
		})
	}
}

func TestAuthenticatedMiddleware(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET_KEY")
	defer func() {
		os.Setenv("JWT_SECRET_KEY", originalSecret)
	}()

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

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		expectedUser   *entity.User
	}{
		{
			name: "Valid authentication",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/user/list", nil)
				req.Header.Set("Authorization", "Bearer "+validTokenString)
				return req
			},
			expectedStatus: http.StatusOK,
			expectedUser: &entity.User{
				Username: "testuser",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedUser *entity.User
			var handlerCalled bool

			// Mock handler that will capture the user from context
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				userVal := r.Context().Value("user")
				if userVal != nil {
					capturedUser = userVal.(*entity.User)
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware chain
			middleware := Authenticated(nextHandler)

			// Create a response recorder and request
			recorder := httptest.NewRecorder()
			req := tt.setupRequest()

			// Call the middleware
			middleware.ServeHTTP(recorder, req)

			// Check results
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.expectedUser != nil {
				assert.True(t, handlerCalled, "Next handler should have been called")
				assert.NotNil(t, capturedUser, "User should be in context")
				assert.Equal(t, tt.expectedUser.Username, capturedUser.Username)
			} else {
				assert.False(t, handlerCalled, "Next handler should not have been called")
			}
		})
	}
}
