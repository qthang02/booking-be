package request

import "github.com/qthang02/booking/types"

type CreateRoomRequest struct {
	RoomNumber uint64 `json:"room_number" binding:"required"`
	CategoryId uint   `json:"category_id" binding:"required"`
}

type UpdateRoomRequest struct {
	RoomNumber uint64           `json:"room_number"`
	Status     types.RoomStatus `json:"status"`
	CategoryId uint             `json:"category_id"`
}
