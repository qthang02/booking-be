package enities

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name     string     `json:"name"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	Birthday *time.Time `json:"birthday"`
	Gender   bool       `json:"gender"`
	Address  string     `json:"address"`
	Password string     `json:"password"`
	Orders   []Order    `json:"orders" gorm:"foreignKey:UserID"`
}
