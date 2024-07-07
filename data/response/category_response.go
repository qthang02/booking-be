package response

import (
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
)

type ListCategoriesResponse struct {
	Categories []*enities.Category `json:"categories"`
	Paging     *requset.Paging     `json:"paging"`
}
