package userrepo

import (
	"context"
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
)

type IUserRepo interface {
	CreateUser(ctx context.Context, req *requset.CreateUserRequest) error
	FindByEmail(ctx context.Context, email string) (*enities.User, error)
	FindByID(ctx context.Context, id int) (*enities.User, error)
	UpdateUser(ctx context.Context, id int, req *requset.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context) ([]*enities.User, error)
}
