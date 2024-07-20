package response

import (
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
)

type ListOrdersResponse struct {
	Orders []*enities.Order
	Paging *request.Paging
}
