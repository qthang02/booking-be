package factory

import (
	"github.com/qthang02/booking/helper"
	"github.com/qthang02/booking/services/biz"
	"github.com/qthang02/booking/services/repo"
)

var (
	userRepo repo.IUserRepo
	userBiz  *biz.UserBiz
)

func Default(config helper.Config) {
	db := helper.ConnectionDB(config)
	userRepo = repo.NewUserRepo(db)

	userBiz = biz.NewUserBiz(userRepo, &config)
}

func GetUserBiz() *biz.UserBiz {
	return userBiz
}
