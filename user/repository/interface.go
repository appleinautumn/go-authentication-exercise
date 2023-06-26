package repository

import (
	"context"

	"imp/assessment/user/entity"
)

type UserRepository interface {
	List(ctx context.Context) ([]*entity.User, error)
	FindOneByUsername(ctx context.Context, username string) (*entity.User, error)
	Create(ctx context.Context, u *entity.User) (*entity.User, error)
}
