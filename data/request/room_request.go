package request

import "github.com/qthang02/booking/types"

type CreateRoomRequest struct {
	RoomNumber uint64 `json:"roomNumber" binding:"required"`
	CategoryId uint   `json:"categoryId" binding:"required"`
}

type UpdateRoomRequest struct {
	RoomNumber uint64           `json:"roomNumber"`
	Status     types.RoomStatus `json:"status"`
	CategoryId uint             `json:"categoryId"`
}
