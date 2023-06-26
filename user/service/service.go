package service

import (
	"context"

	"imp/assessment/user/entity"
	"imp/assessment/user/repository"
)

type userService struct {
	repository repository.UserRepository
}

func NewService(repo repository.UserRepository) UserService {
	return &userService{
		repository: repo,
	}
}

func (s *userService) List(ctx context.Context) ([]*entity.User, error) {
	users, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
