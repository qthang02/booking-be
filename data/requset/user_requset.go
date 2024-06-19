package requset

import "github.com/qthang02/booking/common"

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
