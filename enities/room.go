package enities

import (
	"github.com/qthang02/booking/types"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	RoomNumber uint64           `json:"room_number"`
	Status     types.RoomStatus `json:"status"`
	CategoryId uint             `json:"category_id"`
}
