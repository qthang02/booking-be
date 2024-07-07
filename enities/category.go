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
	AvailableRooms uint64             `json:"available_rooms"`
	Type           types.CategoryType `json:"type"`
	Rooms          []Room             `json:"rooms" gorm:"foreignKey:CategoryID"`
}
