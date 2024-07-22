package userrepo

import (
	"context"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
)

type IUserRepo interface {
	Save(ctx context.Context, user *enities.User) error
	FindByEmail(ctx context.Context, email string) (*enities.User, error)
	FindByID(ctx context.Context, id int) (*enities.User, error)
	UpdateUser(ctx context.Context, id int, req *request.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, paging *request.Paging) ([]*enities.User, error)
}
