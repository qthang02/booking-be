package response

import "github.com/qthang02/booking/data/request"

type PaginatedResponse struct {
	Data   interface{}     `json:"data"`
	Paging *request.Paging `json:"paging"`
}
