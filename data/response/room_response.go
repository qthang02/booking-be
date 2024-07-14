package response

import (
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
)

type ListRoomsResponse struct {
	Rooms  []*enities.Room `json:"rooms"`
	Paging *request.Paging `json:"paging"`
}
