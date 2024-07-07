package userrepo

import (
	"fmt"
	"github.com/jinzhu/copier"
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

func (repo *UserRepo) Save(user *enities.User) error {
	if err := repo.db.Create(user).Error; err != nil {
		log.Error().Err(err).Msg("UserRepo.Save cannot save user")
		return err
	}
	return nil
}

func (repo *UserRepo) CreateUser(req *requset.CreateUserRequest) error {
	var data *enities.User

	err := copier.Copy(data, req)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("UserRepo.CreateUser copier.Copy failed err: %s with req: %s", err, req))
		return err
	}

	resp := repo.db.Create(data)

	return resp.Error
}

func (repo *UserRepo) FindByEmail(email string) (*enities.User, error) {
	var user enities.User

	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Error().Err(err).Msg("UserRepo.FindByEmail User not found")
		return &enities.User{}, err
	}
	return &user, nil
}

func (repo *UserRepo) FindByID(id int) (*enities.User, error) {
	var user enities.User

	err := repo.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		log.Error().Err(err).Msg("UserRepo.FindByID User not found")
		return &enities.User{}, err
	}

	return &user, nil
}

func (repo *UserRepo) UpdateUser(id int, req *requset.UpdateUserRequest) error {

	var user enities.User

	err := copier.Copy(&user, req)
	if err != nil {
		log.Error().Err(err).Msg("copier.Copy User")
		return err
	}

	if err := repo.db.Where(id).Updates(&user).Error; err != nil {
		log.Error().Err(err).Msg("UserRepo.UpdateUser User not found")
		return err
	}

	return nil
}

func (repo *UserRepo) DeleteUser(id int) error {
	if err := repo.db.Where("id = ?", id).Delete(&enities.User{}).Error; err != nil {
		log.Error().Err(err).Msg("UserRepo.DeleteUser User not found")
		return err
	}

	return nil
}

func (repo *UserRepo) GetAllUsers() ([]*enities.User, error) {
	var users []*enities.User
	err := repo.db.Find(&users).Error
	if err != nil {
		log.Error().Err(err).Msg("UserRepo.GetAllUsers User not found")
		return nil, err
	}
	return users, nil
}
