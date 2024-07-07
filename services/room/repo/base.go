package roomrepo

import (
	"context"
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
)

type IRoomRepo interface {
	ListRooms(ctx context.Context, paging *requset.Paging) ([]*enities.Room, error)
	GetRoom(ctx context.Context, id int) (*enities.Room, error)
	CreateRoom(ctx context.Context, request *requset.CreateRoomRequest) error
	UpdateRoom(ctx context.Context, request *requset.UpdateRoomRequest) error
	DeleteRoom(ctx context.Context, id int) error
}
