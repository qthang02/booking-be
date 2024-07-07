package requset

import "github.com/qthang02/booking/types"

type CreateCategoryRequest struct {
	Name           string             `json:"name" binding:"required"`
	Description    string             `json:"description"`
	ImageLink      string             `json:"image_link"`
	Price          float64            `json:"price" binding:"required"`
	AvailableRooms uint64             `json:"available_rooms"`
	Type           types.CategoryType `json:"type" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	ImageLink      string             `json:"image_link"`
	Price          float64            `json:"price"`
	AvailableRooms uint64             `json:"available_rooms"`
	Type           types.CategoryType `json:"type"`
}
