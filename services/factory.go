package services

import (
	authen "github.com/qthang02/booking/services/authen/biz"
	user "github.com/qthang02/booking/services/user/biz"
	userrepo "github.com/qthang02/booking/services/user/repo"
)

var (
	userRepo  userrepo.IUserRepo
	userBiz   *user.UserBiz
	authenBiz *authen.AuthenBiz
)

func Default(config util.Config) {
	db := util.ConnectionDB(config)
	userRepo = userrepo.NewUserRepo(db)

	userBiz = user.NewUserBiz(userRepo, &config)
	authenBiz = authen.NewAuthenBiz(userRepo, &config)
}

func GetUserBiz() *user.UserBiz {
	return userBiz
}

func GetAuthenBiz() *authen.AuthenBiz {
	return authenBiz
}
