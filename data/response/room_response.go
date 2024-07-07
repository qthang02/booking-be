package response

import (
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
)

type ListRoomsResponse struct {
	Rooms  []*enities.Room `json:"rooms"`
	Paging *requset.Paging `json:"paging"`
}
