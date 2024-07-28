package request

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
	Phone    string `json:"phone"`
	Gender   bool   `json:"gender"`
	Address  string `json:"address"`
}
