package requset

import (
	"github.com/qthang02/booking/common"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Username string          `json:"username" binding:"required"`
	Email    string          `json:"email" binding:"required"`
	UserType common.UserType `json:"user_type" binding:"required"`
	Password string          `form:"password" binding:"required"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ListUsersRequest struct {
}

type GetUserRequest struct {
	ID uint `json:"id" binding:"required"`
}

type UpdateUserRequest struct {
	gorm.Model
	Username string          `json:"username"`
	Email    string          `json:"email"`
	UserType common.UserType `json:"user_type"`
	Password string          `json:"password"`
}
