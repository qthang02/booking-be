package requset

import "github.com/qthang02/booking/types"

type CreateRoomRequest struct {
	RoomName   uint64 `json:"room_name" binding:"required"`
	CategoryId uint   `json:"category_id" binding:"required"`
}

type UpdateRoomRequest struct {
	RoomName   uint64           `json:"room_name"`
	Status     types.RoomStatus `json:"status"`
	CategoryId uint             `json:"category_id"`
}
