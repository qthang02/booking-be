package services

import (
	authen "github.com/qthang02/booking/services/authen/biz"
	categotybiz "github.com/qthang02/booking/services/category/biz"
	categoryrepo "github.com/qthang02/booking/services/category/repo"
	roombiz "github.com/qthang02/booking/services/room/biz"
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
	categoryBiz  *categotybiz.CategoryBiz
	roomBiz      *roombiz.RoomBiz
)

func Default(config util.Config) {
	// repo
	db := util.ConnectionDB(config)
	userRepo = userrepo.NewUserRepo(db)
	roomRepo = roomrepo.NewRoomRepo(db)
	categoryRepo = categoryrepo.NewCategoryRepo(db)

	// biz
	userBiz = user.NewUserBiz(userRepo, &config)
	authenBiz = authen.NewAuthenBiz(userRepo, &config)
	categoryBiz = categotybiz.NewCategoryBiz(categoryRepo, &config)
	roomBiz = roombiz.NewRoomBiz(roomRepo, &config)
}

func GetUserBiz() *user.UserBiz {
	return userBiz
}

func GetAuthenBiz() *authen.AuthenBiz {
	return authenBiz
}

func GetCategoryBiz() *categotybiz.CategoryBiz {
	return categoryBiz
}

func GetRoomBiz() *roombiz.RoomBiz {
	return roomBiz
}
