package enities

import (
	"github.com/qthang02/booking/types"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	ImageLink      string             `json:"image_link"`
	Price          float64            `json:"price"`
	AvailableRooms int64              `json:"available_rooms" gorm:"-"`
	Type           types.CategoryType `json:"type"`
	Rooms          []Room             `json:"rooms" gorm:"foreignKey:CategoryId"`
}
