package service

import (
	"context"
	"go-authentication-exercise/user/entity"
)

type UserService interface {
	List(ctx context.Context) ([]*entity.User, error)
	Count(ctx context.Context) (int, error)
}
