package roomrepo

import (
	"github.com/qthang02/booking/enities"
	"gorm.io/gorm"
)

type RoomRepo struct {
	db *gorm.DB
}

func NewRoomRepo(db *gorm.DB) IRoomRepo {
	db.AutoMigrate(&enities.Room{})

	return &RoomRepo{
		db: db,
	}
}
