package userrepo

import (
	"context"
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
	err := db.AutoMigrate(&enities.User{})
	if err != nil {
		log.Error().Msgf("failed to auto migrate user database")
		return nil
	}

	return &UserRepo{db: db}
}

func (repo *UserRepo) CreateUser(_ context.Context, req *requset.CreateUserRequest) error {
	var data *enities.User

	err := copier.Copy(data, req)
	if err != nil {
		log.Error().Err(err).Msgf("UserRepo.CreateUser copier.Copy failed err: %s with req: %v", err, req)
		return err
	}

	resp := repo.db.Create(data)

	return resp.Error
}

func (repo *UserRepo) FindByEmail(_ context.Context, email string) (*enities.User, error) {
	var user enities.User

	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Error().Err(err).Msgf("UserRepo.FindByEmail User not found err: %s with email: %v", err, email)
		return &enities.User{}, err
	}
	return &user, nil
}

func (repo *UserRepo) FindByID(_ context.Context, id int) (*enities.User, error) {
	var user enities.User

	err := repo.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		log.Error().Err(err).Msgf("UserRepo.FindByID User not found err: %s with id: %v", err, id)
		return &enities.User{}, err
	}

	return &user, nil
}

func (repo *UserRepo) UpdateUser(_ context.Context, id int, req *requset.UpdateUserRequest) error {

	var user enities.User

	err := copier.Copy(&user, req)
	if err != nil {
		log.Error().Err(err).Msgf("UserRepo.UpdateUser cannot copy user requset err: %s with request: %v", err, req)
		return err
	}

	if err := repo.db.Where(id).Updates(&user).Error; err != nil {
		log.Error().Err(err).Msgf("UserRepo.UpdateUser User not found err: %s with id: %v", err, id)
		return err
	}

	return nil
}

func (repo *UserRepo) DeleteUser(_ context.Context, id int) error {
	if err := repo.db.Where("id = ?", id).Delete(&enities.User{}).Error; err != nil {
		log.Error().Err(err).Msgf("UserRepo.DeleteUser cannot delete user requset err: %s with id: %v", err, id)
		return err
	}

	return nil
}

func (repo *UserRepo) ListUsers(_ context.Context) ([]*enities.User, error) {
	var users []*enities.User
	err := repo.db.Find(&users).Error
	if err != nil {
		log.Error().Err(err).Msgf("UserRepo.GetAllUsers cannot list users err: %s", err)
		return nil, err
	}
	return users, nil
}
