package service

import (
	"context"

	"imp/assessment/user/entity"
)

type AuthService interface {
	Login(ctx context.Context, username string, password string) (string, error)
	Signup(ctx context.Context, username string, fullname string, password string) (*entity.User, error)
}
