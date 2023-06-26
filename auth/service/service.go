package service

import (
	"context"
	"errors"

	"imp/assessment/user/entity"
	"imp/assessment/user/repository"

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
	return "", nil
}

// Function to hash the user's password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
