package repo

import (
	"github.com/qthang02/booking/data/requset"
)

type IUserRepo interface {
	CreateUser(req *requset.CreateUserRequest) error
}
