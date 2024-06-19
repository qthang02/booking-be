package biz

import (
	"encoding/json"
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/helper"
	"github.com/qthang02/booking/services/repo"
	"github.com/rs/zerolog/log"
	"net/http"
)

type UserBiz struct {
	userRepo repo.IUserRepo
}

func NewUserBiz(userRepo repo.IUserRepo) *UserBiz {
	return &UserBiz{userRepo: userRepo}
}

func (biz *UserBiz) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user requset.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser cannot hash password")
		return
	}

	user.Password = hashPassword

	err = biz.userRepo.CreateUser(&user)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser cannot create user")
		return
	}
}

// TODO: list user

// TODO: get user

// TODO: delete user

// TODO: update user

func (biz *UserBiz) Login(w http.ResponseWriter, r *http.Request) {
	log.Log().Msg("UserBiz.Login request")
	var login requset.LoginUserRequest

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Msg("UserBiz.Login cannot decode request")
		return
	}

	user, err := biz.userRepo.FindByUsername(login.Username)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.Login cannot find user")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = helper.VerifyPassword(user.Password, login.Password)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.Login cannot verify password")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	config, _ := helper.LoadConfig(".")

	token, err := helper.GenerateToken(config.TokenExpiresIn, user.ID, config.TokenSecret)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.Login cannot generate token")
		return
	}

	w.Write([]byte(token))
}
