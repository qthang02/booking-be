package cmd

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qthang02/booking/services"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

func setupHttpRoutes(server *echo.Echo) {

	api := server.Group("/api/v1")
	{
		api.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "I'm still alive")
		})

		user := api.Group("/user")
		{
			user.GET("/:id", services.GetUserBiz().GetUserById)
			user.PUT("/:id", services.GetUserBiz().UpdateUser)
			user.GET("", services.GetUserBiz().ListUsers)
			user.DELETE("/:id", services.GetUserBiz().DeleteUserById)
			user.POST("", services.GetUserBiz().CreateUser)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/register", services.GetAuthenBiz().RegisterUser)
			auth.POST("/login", services.GetAuthenBiz().Login)
		}

		category := api.Group("/category")
		{
			category.GET("/:id", services.GetCategoryBiz().Get)
			category.GET("", services.GetCategoryBiz().List)
			category.POST("", services.GetCategoryBiz().Create)
			category.PUT("/:id", services.GetCategoryBiz().Update)
			category.DELETE("/:id", services.GetCategoryBiz().Delete)
		}

		room := api.Group("/room")
		{
			room.GET("/:id", services.GetRoomBiz().Get)
			room.GET("", services.GetRoomBiz().List)
			room.POST("", services.GetRoomBiz().Create)
			room.PUT("/:id", services.GetRoomBiz().Update)
			room.DELETE("/:id", services.GetRoomBiz().Delete)
		}

		order := api.Group("/order")
		{
			order.GET("/:id", services.GetOrderBiz().GetOrder)
			order.GET("", services.GetOrderBiz().ListOrders)
			order.POST("", services.GetOrderBiz().CreateOrder)
			order.PUT("/:id", services.GetOrderBiz().UpdateOrder)
			order.DELETE("/:id", services.GetOrderBiz().DeleteOrder)
		}
	}
}

func Run() {
	server := echo.New()
	server.Use(middleware.CORS())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	logger := zerolog.New(os.Stdout)
	config, err := util.LoadConfig(logger)
	if err != nil {
		logger.Error().Err(err).Msg("cannot load config")
	}
	services.Default(config)
	setupHttpRoutes(server)

	PORT := "8080"
	server.Logger.Fatal(server.Start(fmt.Sprintf(":%s", PORT)))
}
