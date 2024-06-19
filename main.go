package main

import (
	"github.com/qthang02/booking/factory"
	"github.com/qthang02/booking/handlers"
	"github.com/qthang02/booking/helper"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Error().Err(err).Msg("cannot load config")
	}
	factory.Default(config)
	handlers.InitUserHandler(router)

	server := http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	}

	server.ListenAndServe()
}
