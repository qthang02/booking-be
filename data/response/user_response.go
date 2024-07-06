package response

import (
	"time"
)

type UserDTOResponse struct {
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Name     string     `json:"name"`
	Phone    string     `json:"phone"`
	Birthday *time.Time `json:"birthday"`
	Gender   bool       `json:"gender"`
	Address  string     `json:"address"`
}
