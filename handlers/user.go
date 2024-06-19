package handlers

import (
	"github.com/qthang02/booking/services/factory"
	"net/http"
)

func InitUserHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/user", factory.GetUserBiz().CreateUser)
	mux.HandleFunc("POST /api/v1/user/login", factory.GetUserBiz().Login)
}
