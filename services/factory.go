package services

import (
	authen "github.com/qthang02/booking/services/authen/biz"
	categoryrepo "github.com/qthang02/booking/services/category/repo"
	"github.com/qthang02/booking/services/room/repo"
	user "github.com/qthang02/booking/services/user/biz"
	userrepo "github.com/qthang02/booking/services/user/repo"
	"github.com/qthang02/booking/util"
)

var (
	userRepo     userrepo.IUserRepo
	roomRepo     roomrepo.IRoomRepo
	categoryRepo categoryrepo.ICategoryRepo
	userBiz      *user.UserBiz
	authenBiz    *authen.AuthenBiz
)

func Default(config util.Config) {
	db := util.ConnectionDB(config)
	userRepo = userrepo.NewUserRepo(db)
	roomRepo = roomrepo.NewRoomRepo(db)
	categoryRepo = categoryrepo.NewCategoryRepo(db)

	userBiz = user.NewUserBiz(userRepo, &config)
	authenBiz = authen.NewAuthenBiz(userRepo, &config)
}

func GetUserBiz() *user.UserBiz {
	return userBiz
}

func GetAuthenBiz() *authen.AuthenBiz {
	return authenBiz
}
