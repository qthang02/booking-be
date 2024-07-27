package cmd

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	middlewarecustom "github.com/qthang02/booking/middleware"
	"github.com/qthang02/booking/services"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

func setupHttpRoutes(server *echo.Echo, config util.Config) {

	api := server.Group("/api/v1")
	{
		api.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "I'm still alive")
		})

		user := api.Group("/user", middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
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
			auth.GET("/profile", services.GetAuthenBiz().Profile, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Customer, util.Staff}))
		}

		category := api.Group("/category")
		{
			category.GET("/:id", services.GetCategoryBiz().Get)
			category.GET("", services.GetCategoryBiz().List)
			category.POST("", services.GetCategoryBiz().Create, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
			category.PUT("/:id", services.GetCategoryBiz().Update, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
			category.DELETE("/:id", services.GetCategoryBiz().Delete, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
		}

		room := api.Group("/room")
		{
			room.GET("/:id", services.GetRoomBiz().Get)
			room.GET("", services.GetRoomBiz().List)
			room.POST("", services.GetRoomBiz().Create, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
			room.PUT("/:id", services.GetRoomBiz().Update, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
			room.DELETE("/:id", services.GetRoomBiz().Delete, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
		}

		order := api.Group("/order")
		{
			order.GET("/:id", services.GetOrderBiz().GetOrder)
			order.GET("", services.GetOrderBiz().ListOrders)
			order.POST("", services.GetOrderBiz().CreateOrder, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
			order.PUT("/:id", services.GetOrderBiz().UpdateOrder, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
			order.DELETE("/:id", services.GetOrderBiz().DeleteOrder, middlewarecustom.JWTAuth(config.TokenSecret, []string{util.Admin, util.Staff}))
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
	setupHttpRoutes(server, config)

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s", config.ServerAddress)))
}
