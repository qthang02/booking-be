package util

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func ConnectionDB(config Config) *gorm.DB {
	dsn := config.DBSource
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Info().Msg("Successfully connected to database")
			return db
		}
		log.Error().Err(err).Msgf("Failed to connect to database (attempt %d/5)", i+1)
		time.Sleep(5 * time.Second)
	}
	log.Fatal().Err(err).Msg("Failed to connect to database after 5 attempts")
	return nil
}
