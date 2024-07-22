package userrepo

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/database"
	"github.com/qthang02/booking/enities"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	userRepo *UserRepo
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	userRepo = &UserRepo{}

	err := db.AutoMigrate(&enities.User{})
	if err != nil {
		log.Error().Msgf("failed to auto migrate user database")
		return nil
	}

	userRepo.db = db

	err = userRepo.initUserDB()
	if err != nil {
		log.Error().Msgf("failed to init user database")
		return nil
	}

	return userRepo
}

func (repo *UserRepo) initUserDB() error {
	users, err := repo.ListUsers(context.Background(), &request.Paging{})
	if err != nil {
		log.Error().Msgf("initUserDB: failed to list users")
		return err
	}

	if len(users) == 0 {
		for _, user := range database.InitUsersDataDefault() {
			user.Password, err = util.HashPassword(user.Password)
			if err != nil {
				log.Error().Msgf("initUserDB: failed to hash password")
				return err
			}
			err := repo.Save(context.Background(), user)
			if err != nil {
				log.Error().Msgf("initUserDB: failed to save user")
				return err
			}
		}
	}

	return nil
}

func (repo *UserRepo) Save(_ context.Context, user *enities.User) error {
	if err := repo.db.Create(user).Error; err != nil {
		log.Error().Err(err).Msg("UserRepo.Save cannot save user")
		return err
	}
	return nil
}

func (repo *UserRepo) FindByEmail(ctx context.Context, email string) (*enities.User, error) {
	var user enities.User

	err := repo.db.WithContext(ctx).
		Preload("Orders").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msgf("UserRepo.FindByEmail: User not found with email: %v", email)
			return nil, err
		}
		log.Error().Err(err).Msgf("UserRepo.FindByEmail: Error fetching user with email: %v", email)
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepo) FindByID(ctx context.Context, id int) (*enities.User, error) {
	var user enities.User

	err := repo.db.WithContext(ctx).
		Preload("Orders").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msgf("UserRepo.FindByID: User not found with id: %v", id)
			return nil, err
		}
		log.Error().Err(err).Msgf("UserRepo.FindByID: Error fetching user with id: %v", id)
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepo) UpdateUser(_ context.Context, id int, req *request.UpdateUserRequest) error {

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

func (repo *UserRepo) ListUsers(ctx context.Context, paging *request.Paging) ([]*enities.User, error) {
	var users []*enities.User

	offset := (paging.Page - 1) * paging.Limit

	result := repo.db.
		Preload("Orders").
		Limit(paging.Limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("UserRepo.ListUsers: failed to list users")
		return nil, result.Error
	}

	var totalCount int64
	if err := repo.db.Model(&enities.User{}).Count(&totalCount).Error; err != nil {
		log.Error().Err(err).Msg("UserRepo.ListUsers: failed to count total users")
		return nil, err
	}
	paging.Total = totalCount

	return users, nil
}
