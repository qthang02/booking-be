package response

import (
	"gorm.io/gorm"
	"time"
)

type UserDTOResponse struct {
	gorm.Model
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Name     string     `json:"name"`
	Phone    string     `json:"phone"`
	Birthday *time.Time `json:"birthday"`
	Gender   bool       `json:"gender"`
	Address  string     `json:"address"`
}
