package response

import (
	"time"
)

type UserDTOResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Birthday  *time.Time        `json:"birthday,omitempty"`
	Gender    bool              `json:"gender"`
	Address   string            `json:"address"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Orders    []OrderSummaryDTO `json:"orders,omitempty"`
	Role      string            `json:"role,omitempty"`
}
