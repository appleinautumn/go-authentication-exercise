package service

import (
	"context"
	"errors"
	"os"
	"time"

	"go-authentication-exercise/user/entity"
	"go-authentication-exercise/user/repository"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repository repository.UserRepository
}

func NewService(repo repository.UserRepository) AuthService {
	return &authService{
		repository: repo,
	}
}

func (s *authService) Signup(ctx context.Context, username string, fullname string, password string) (*entity.User, error) {
	// get user
	user, err := s.repository.FindOneByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("username exists")
	}

	// hash the password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		Id:       uuid.New(),
		Username: username,
		Fullname: fullname,
		Password: hashedPassword,
	}

	res, err := s.repository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *authService) Login(ctx context.Context, username string, password string) (string, error) {
	// get user
	user, err := s.repository.FindOneByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user doesn't exist")
	}

	// Simulate user login
	loginSuccess := checkPasswordHash(password, user.Password)
	if !loginSuccess {
		return "", errors.New("invalid login")
	}

	// success, now generate the token
	accessToken, err := generateJwtAccessToken(user.Username)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Function to hash the user's password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// Function to check if the provided password matches the hashed password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJwtAccessToken(username string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
