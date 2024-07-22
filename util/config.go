package util

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBSource       string
	ServerAddress  string
	TokenSecret    string
	TokenExpiresIn time.Duration
	TokenMaxAge    int
}

func LoadConfig(logger zerolog.Logger) (config Config, err error) {
	logger.Info().Timestamp().Msg("Loading config...")

	config.DBSource = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))
	config.ServerAddress = os.Getenv("SERVER_ADDRESS")
	config.TokenSecret = os.Getenv("TOKEN_SECRET")
	config.TokenExpiresIn, _ = time.ParseDuration(os.Getenv("TOKEN_EXPIRED_IN"))
	config.TokenMaxAge, _ = strconv.Atoi(os.Getenv("TOKEN_MAXAGE"))

	return config, nil
}
