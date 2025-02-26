package repository

import (
	"context"

	"go-authentication-exercise/user/entity"
)

type UserRepository interface {
	List(ctx context.Context) ([]*entity.User, error)
	Count(ctx context.Context) (int, error)
	FindOneByUsername(ctx context.Context, username string) (*entity.User, error)
	Create(ctx context.Context, u *entity.User) (*entity.User, error)
}
