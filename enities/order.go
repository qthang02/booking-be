package enities

import (
	"github.com/qthang02/booking/types"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	GuestNumber  uint64             `json:"guestNumber"`
	Price        float64            `json:"price"`
	Description  string             `json:"description"`
	Checkin      time.Time          `json:"checkin"`
	Checkout     time.Time          `json:"checkout"`
	CategoryType types.CategoryType `json:"categoryType"`
	RoomNumber   uint64             `json:"roomNumber"`
	UserID       uint               `json:"userID"`
}
