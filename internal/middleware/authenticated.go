package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go-authentication-exercise/internal/user/entity"
	"go-authentication-exercise/internal/util"

	"github.com/golang-jwt/jwt"
)

// Authenticated middleware checks if the request has a valid JWT token
// and adds the user information to the request context
func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from header or query parameter
		token, err := extractToken(r)
		if err != nil {
			util.Error(w, http.StatusUnauthorized, nil, err.Error())
			return
		}

		// Validate the token
		claims, err := validateToken(token)
		if err != nil {
			util.Error(w, http.StatusUnauthorized, nil, "Invalid token: "+err.Error())
			return
		}

		// Get username from token claims
		username, err := getUsernameFromJwt(claims)
		if err != nil {
			util.Error(w, http.StatusUnauthorized, nil, "Invalid token: "+err.Error())
			return
		}

		// Create user context and proceed with request
		user := &entity.User{
			Username: username,
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken gets the JWT token from either the Authorization header
// or from a query parameter
func extractToken(r *http.Request) (string, error) {
	// Check Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Fall back to query parameter
		authHeader = r.URL.Query().Get("authorization")
	}

	if authHeader == "" {
		return "", errors.New("authorization token is required")
	}

	// Check for Bearer token format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

// validateToken verifies that the token is valid and returns its claims
func validateToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET_KEY is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// getUsernameFromJwt extracts the username from JWT claims
func getUsernameFromJwt(data map[string]interface{}) (string, error) {
	if val, ok := data["username"]; ok {
		username, ok := val.(string)
		if !ok || username == "" {
			return "", errors.New("invalid username in token")
		}
		return username, nil
	}

	return "", errors.New("username claim not found in token")
}
