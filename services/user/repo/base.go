package userrepo

import (
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
)

type IUserRepo interface {
	Save(user *enities.User) error
	CreateUser(req *requset.CreateUserRequest) error
	FindByEmail(email string) (*enities.User, error)
	FindByID(id int) (*enities.User, error)
	UpdateUser(id int, req *requset.UpdateUserRequest) error
	DeleteUser(id int) error
	GetAllUsers() ([]*enities.User, error)
}
