package cmd

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qthang02/booking/services"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog"
	"os"
)

func setupHttpRoutes(server *echo.Echo) {

	api := server.Group("/api/v1")
	{
		user := api.Group("/user")
		{
			user.GET("/:id", services.GetUserBiz().GetUserById)
			user.PUT("/:id", services.GetUserBiz().UpdateUser)
			user.GET("/users", services.GetUserBiz().ListUsers)
			user.DELETE("/:id", services.GetUserBiz().DeleteUserById)
			user.POST("/user", services.GetUserBiz().CreateUser)

		}

		auth := api.Group("/auth")
		{
			auth.POST("/register", services.GetAuthenBiz().RegisterUser)
			auth.POST("/login", services.GetAuthenBiz().Login)
		}
	}
}

func Run() {
	server := echo.New()
	server.Use(middleware.CORS())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	logger := zerolog.New(os.Stdout)
	config, err := util.LoadConfig(".", logger)
	if err != nil {
		logger.Error().Err(err).Msg("cannot load config")
	}
	services.Default(config)
	setupHttpRoutes(server)

	server.Logger.Fatal(server.Start(config.ServerAddress))
}
