package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go-authentication-exercise/user/entity"
	"go-authentication-exercise/util"

	"github.com/golang-jwt/jwt"
)

func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")

		if authorizationHeader == "" {
			authorizationHeader = r.URL.Query().Get("authorization")
		}

		if authorizationHeader == "" {
			util.Error(w, http.StatusUnauthorized, nil, "An authorization header is required")
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")

		if len(bearerToken) != 2 {
			util.Error(w, http.StatusUnauthorized, nil, "An authorization header is required")
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error")
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			util.Error(w, http.StatusUnauthorized, nil, "Invalid token: "+err.Error())
			return
		}

		if !token.Valid {
			util.Error(w, http.StatusUnauthorized, nil, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			util.Error(w, http.StatusUnauthorized, nil, "Invalid token: "+err.Error())
			return
		}

		username, err := getUsernameFromJwt(claims)

		if err != nil {
			util.Error(w, http.StatusUnauthorized, nil, "Invalid token: "+err.Error())
			return
		}

		user := &entity.User{
			Username: username,
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getUsernameFromJwt(data map[string]interface{}) (string, error) {
	if val, ok := data["username"]; ok {
		username, _ := val.(string)
		return username, nil
	}

	return "", errors.New("invalid Token")
}
