package service

import (
	"context"
	"imp/assessment/user/entity"
)

type UserService interface {
	List(ctx context.Context) ([]*entity.User, error)
}
