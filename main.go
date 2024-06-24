package main

import (
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/handlers"
	"github.com/qthang02/booking/helper"
	"github.com/qthang02/booking/services/factory"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	server := echo.New()
	logger := zerolog.New(os.Stdout)
	config, err := helper.LoadConfig(".", logger)
	if err != nil {
		logger.Error().Err(err).Msg("cannot load config")
	}
	factory.Default(config)
	handlers.InitUserHandler(server)

	server.Logger.Fatal(server.Start(config.ServerAddress))
}
