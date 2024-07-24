package request

import (
	"time"
)

type CreateUserRequest struct {
	Username string     `json:"username" binding:"required"`
	Email    string     `json:"email" binding:"required"`
	Phone    string     `json:"phone"`
	Birthday *time.Time `json:"birthday"`
	Gender   bool       `json:"gender"`
	Address  string     `json:"address"`
	Password string     `json:"password" binding:"required"`
}

type ListUsersRequest struct {
}

type GetUserRequest struct {
	ID uint `json:"id" binding:"required"`
}

type UpdateUserRequest struct {
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	Birthday *time.Time `json:"birthday"`
	Gender   bool       `json:"gender"`
	Address  string     `json:"address"`
	Password string     `json:"password"`
}
