package helper

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DBSource       string        `mapstructure:"DB_SOURCE"`
	ServerAddress  string        `mapstructure:"SERVER_ADDRESS"`
	TokenSecret    string        `mapstructure:"TOKEN_SECRET"`
	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRED_IN"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAXAGE"`
}

func LoadConfig(path string) (config Config, err error) {
	log.Log().Msg("Loading config...")

	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msg("Error reading config")
		return
	}

	err = viper.Unmarshal(&config)
	return
}
