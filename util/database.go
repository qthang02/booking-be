package util

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectionDB(config Config) *gorm.DB {
	dsn := config.DBSource
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("NewDatabase cannot connect database error")
	}

	return db
}
