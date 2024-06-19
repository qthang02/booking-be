package handlers

import (
	"github.com/qthang02/booking/factory"
	"net/http"
)

func InitUserHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/user", factory.GetUserBiz().CreateUser)
}
