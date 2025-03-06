package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
