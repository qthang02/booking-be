package response

import (
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
)

type ListCategoriesResponse struct {
	Categories []*enities.Category `json:"categories"`
	Paging     *request.Paging     `json:"paging"`
}
