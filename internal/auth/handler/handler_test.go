package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the AuthService interface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, username string, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        map[string]interface{}
		setupMock          func(*MockAuthService)
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name: "Successful login",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
			},
			setupMock: func(mockService *MockAuthService) {
				mockService.On("Login", mock.Anything, "testuser", "password123").
					Return("jwt-token-here", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "",
				"data":    "jwt-token-here",
			},
		},
		{
			name: "Invalid credentials",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "wrongpassword",
			},
			setupMock: func(mockService *MockAuthService) {
				mockService.On("Login", mock.Anything, "testuser", "wrongpassword").
					Return("", errors.New("invalid credentials"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "invalid credentials",
				"data":    nil,
			},
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"username": "u", // Too short, validation should fail
				"password": "pwd",
			},
			setupMock: func(mockService *MockAuthService) {
				// Service mock should not be called since validation fails
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "Validation error",
				// We're not checking the exact validation errors here
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(MockAuthService)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			// Create the handler with mock service
			handler := NewAuthHandler(mockService)

			// Create a request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			res := httptest.NewRecorder()

			// Call the handler
			handler.Login(res, req)

			// Check status code
			assert.Equal(t, tt.expectedStatusCode, res.Code)

			// Parse response
			var responseBody map[string]interface{}
			err := json.Unmarshal(res.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			// Check response message
			assert.Equal(t, tt.expectedResponse["message"], responseBody["message"])

			// For successful responses, check data
			if tt.expectedStatusCode == http.StatusOK {
				assert.Equal(t, tt.expectedResponse["data"], responseBody["data"])
			}

			// Verify that all expected mock calls were made
			mockService.AssertExpectations(t)
		})
	}
}
