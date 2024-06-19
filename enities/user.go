package enities

import (
	"github.com/qthang02/booking/common"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string          `json:"name"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	UserType common.UserType `json:"user_type"`
	Status   common.Status   `json:"status"`
	Password string          `json:"password"`
}
