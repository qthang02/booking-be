package request

import (
	"github.com/qthang02/booking/types"
	"time"
)

type CreateOrderRequest struct {
	GuestNumber  uint64             `json:"guestNumber"`
	Price        float64            `json:"price"`
	Description  string             `json:"description"`
	Checkin      time.Time          `json:"checkin"`
	Checkout     time.Time          `json:"checkout"`
	CategoryType types.CategoryType `json:"categoryType"`
	RoomNumber   uint64             `json:"roomNumber"`
	UserID       uint               `json:"userID"`
}

type UpdateOrderRequest struct {
	GuestNumber  uint64             `json:"guestNumber"`
	Price        float64            `json:"price"`
	Description  string             `json:"description"`
	Checkin      time.Time          `json:"checkin"`
	Checkout     time.Time          `json:"checkout"`
	CategoryType types.CategoryType `json:"categoryType"`
	RoomNumber   uint64             `json:"roomNumber"`
}
