package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/services/factory"
)

func InitUserHandler(server *echo.Echo) {
	server.POST("/api/v1/user", factory.GetUserBiz().CreateUser)
	server.GET("/api/v1/user/:id", factory.GetUserBiz().GetUserById)
	server.POST("/api/v1/user/login", factory.GetUserBiz().Login)
	server.PUT("/api/v1/user/:id", factory.GetUserBiz().UpdateUser)
	server.DELETE("/api/v1/user/:id", factory.GetUserBiz().DeleteUserById)
	server.GET("/api/v1/users", factory.GetUserBiz().ListUsers)
}
