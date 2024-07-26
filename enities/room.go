package enities

import (
	"github.com/qthang02/booking/types"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	RoomNumber uint64           `json:"roomNumber" gorm:"uniqueIndex"`
	Status     types.RoomStatus `json:"status"`
	CategoryId int              `json:"categoryId" gorm:"column:category_id"`
	Category   Category         `json:"category" gorm:"foreignKey:CategoryId"`
	Order      Order            `json:"order" gorm:"foreignKey:RoomNumber;references:RoomNumber"`
}
