package response

import (
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
	"time"
)

type ListOrdersResponse struct {
	Orders []*enities.Order
	Paging *request.Paging
}

type OrderSummaryDTO struct {
	ID          uint      `json:"id"`
	GuestNumber uint64    `json:"guest_number"`
	Price       float64   `json:"price"`
	Checkin     time.Time `json:"checkin"`
	Checkout    time.Time `json:"checkout"`
	RoomNumber  uint64    `json:"room_number"`
}
