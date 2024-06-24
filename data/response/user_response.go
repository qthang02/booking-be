package response

import (
	"github.com/qthang02/booking/common"
	"gorm.io/gorm"
)

type UserDTOResponse struct {
	gorm.Model
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Name     string          `json:"name"`
	UserType common.UserType `json:"user_type"`
	Status   common.Status   `json:"status"`
}
