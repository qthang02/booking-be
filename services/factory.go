package services

import (
	authen "github.com/qthang02/booking/services/authen/biz"
	categotybiz "github.com/qthang02/booking/services/category/biz"
	categoryrepo "github.com/qthang02/booking/services/category/repo"
	employeebiz "github.com/qthang02/booking/services/employee/biz"
	employeerepo "github.com/qthang02/booking/services/employee/repo"
	"github.com/qthang02/booking/services/order/biz"
	orderrepo "github.com/qthang02/booking/services/order/repo"
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
	orderRepo    orderrepo.IOrderRepo
	employeeRepo employeerepo.IEmployeeRepo
	userBiz      *user.UserBiz
	authenBiz    *authen.AuthenBiz
	categoryBiz  *categotybiz.CategoryBiz
	roomBiz      *roombiz.RoomBiz
	orderBiz     *orderbiz.OrderBiz
	employeeBiz  *employeebiz.EmployeeBiz
)

func Default(config util.Config) {
	// repo
	db := util.ConnectionDB(config)
	userRepo = userrepo.NewUserRepo(db)
	categoryRepo = categoryrepo.NewCategoryRepo(db)
	roomRepo = roomrepo.NewRoomRepo(db)
	orderRepo = orderrepo.NewOrderRepo(db)
	employeeRepo = employeerepo.NewEmployeeRepo(db)

	// biz
	userBiz = user.NewUserBiz(userRepo, &config)
	authenBiz = authen.NewAuthenBiz(userRepo, &config)
	categoryBiz = categotybiz.NewCategoryBiz(categoryRepo, &config)
	roomBiz = roombiz.NewRoomBiz(roomRepo, &config)
	orderBiz = orderbiz.NewOrderBiz(orderRepo, &config)
	employeeBiz = employeebiz.NewEmployeeBiz(employeeRepo, &config)
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

func GetOrderBiz() *orderbiz.OrderBiz {
	return orderBiz
}

func GetEmployeeBiz() *employeebiz.EmployeeBiz {
	return employeeBiz
}
