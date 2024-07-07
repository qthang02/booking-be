package enities

import (
	"github.com/qthang02/booking/types"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	RoomName   uint64           `json:"room_name"`
	Status     types.RoomStatus `json:"status"`
	CategoryId uint             `json:"category_id"`
}
