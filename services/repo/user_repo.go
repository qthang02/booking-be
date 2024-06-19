package repo

import (
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	db.AutoMigrate(&enities.User{})

	return &UserRepo{db: db}
}

func (repo *UserRepo) CreateUser(req *requset.CreateUserRequest) error {
	data := &enities.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		UserType: req.UserType,
	}

	resp := repo.db.Create(data)

	return resp.Error
}

func (repo *UserRepo) FindByUsername(username string) (enities.User, error) {
	var user enities.User

	err := repo.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Error().Err(err).Msg("UserRepo.FindByUsername User not found")
		return enities.User{}, err
	}
	return user, nil
}

func (repo *UserRepo) FindByID(id int) (enities.User, error) {
	var user enities.User

	err := repo.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		log.Error().Err(err).Msg("UserRepo.FindByID User not found")
		return enities.User{}, err
	}

	return user, nil
}
