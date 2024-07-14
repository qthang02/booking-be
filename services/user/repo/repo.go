package userrepo

import (
	"context"
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
	users, err := repo.ListUsers(context.Background())
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
	if err := repo.db.Save(user).Error; err != nil {
		log.Error().Err(err).Msg("UserRepo.Save cannot save user")
		return err
	}
	return nil
}

func (repo *UserRepo) CreateUser(_ context.Context, req *request.CreateUserRequest) error {
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

func (repo *UserRepo) ListUsers(_ context.Context) ([]*enities.User, error) {
	var users []*enities.User
	err := repo.db.Find(&users).Error
	if err != nil {
		log.Error().Err(err).Msgf("UserRepo.GetAllUsers cannot list users err: %s", err)
		return nil, err
	}
	return users, nil
}
