package repo

import (
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
)

type IUserRepo interface {
	CreateUser(req *requset.CreateUserRequest) error
	FindByUsername(username string) (enities.User, error)
	FindByID(id int) (enities.User, error)
}
