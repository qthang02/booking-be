package repo

import (
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	db.AutoMigrate(&enities.User{})

	return &UserRepo{db: db}
}

func (repo UserRepo) CreateUser(req *requset.CreateUserRequest) error {
	data := &enities.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		UserType: req.UserType,
	}

	resp := repo.db.Create(data)

	return resp.Error
}
