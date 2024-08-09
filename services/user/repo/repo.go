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

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	repo := &UserRepo{db: db}

	err := db.AutoMigrate(&enities.User{})
	if err != nil {
		log.Error().Err(err).Msg("failed to auto migrate user database")
		return nil
	}

	err = repo.initUserDB()
	if err != nil {
		log.Error().Err(err).Msg("failed to init user database")
		return nil
	}

	return repo
}

func (repo *UserRepo) initUserDB() error {
	var count int64
	if err := repo.db.Model(&enities.User{}).Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("initUserDB: failed to count users")
		return err
	}

	log.Info().Msgf("initUserDB: found %d existing users", count)

	if count == 0 {
		for _, user := range database.InitUsersDataDefault() {
			hashedPassword, err := util.HashPassword(user.Password)
			if err != nil {
				log.Error().Err(err).Msg("initUserDB: failed to hash password")
				return err
			}
			user.Password = hashedPassword

			if err := repo.Save(context.Background(), user); err != nil {
				log.Error().Err(err).Msg("initUserDB: failed to save user")
				return err
			}
		}
		log.Info().Msg("initUserDB: default users created")
	}

	return nil
}

func (repo *UserRepo) Save(ctx context.Context, user *enities.User) error {
	if err := repo.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Error().Err(err).Msg("UserRepo.Save: cannot save user")
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
			log.Info().Msgf("UserRepo.FindByEmail: User not found with email: %v", email)
			return nil, gorm.ErrRecordNotFound
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
			log.Info().Msgf("UserRepo.FindByID: User not found with id: %v", id)
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msgf("UserRepo.FindByID: Error fetching user with id: %v", id)
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepo) UpdateUser(ctx context.Context, id int, req *request.UpdateUserRequest) error {
	var user enities.User

	if err := copier.Copy(&user, req); err != nil {
		log.Error().Err(err).Msgf("UserRepo.UpdateUser: cannot copy user request")
		return err
	}

	result := repo.db.WithContext(ctx).Model(&enities.User{}).Where("id = ?", id).Updates(&user)
	if result.Error != nil {
		log.Error().Err(result.Error).Msgf("UserRepo.UpdateUser: cannot update user")
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Info().Msgf("UserRepo.UpdateUser: User not found with id: %v", id)
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo *UserRepo) DeleteUser(ctx context.Context, id int) error {
	result := repo.db.WithContext(ctx).Delete(&enities.User{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msgf("UserRepo.DeleteUser: cannot delete user")
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Info().Msgf("UserRepo.DeleteUser: User not found with id: %v", id)
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo *UserRepo) ListUsers(ctx context.Context, paging *request.Paging) ([]*enities.User, error) {
	var users []*enities.User

	offset := (paging.Page - 1) * paging.Limit

	var totalCount int64
	if err := repo.db.Model(&enities.User{}).Count(&totalCount).Error; err != nil {
		log.Error().Err(err).Msg("UserRepo.ListUsers: failed to count total users")
		return nil, err
	}
	paging.Total = totalCount

	result := repo.db.WithContext(ctx).
		Preload("Orders").
		Where("role = ?", "Customer").
		Limit(paging.Limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("UserRepo.ListUsers: failed to list users")
		return nil, result.Error
	}

	return users, nil
}
