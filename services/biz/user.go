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
