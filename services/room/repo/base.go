package roomrepo

import (
	"context"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
)

type IRoomRepo interface {
	Save(ctx context.Context, room *enities.Room) error
	ListRooms(ctx context.Context, paging *request.Paging) ([]*enities.Room, error)
	GetRoom(ctx context.Context, id int) (*enities.Room, error)
	CreateRoom(ctx context.Context, request *request.CreateRoomRequest) error
	UpdateRoom(ctx context.Context, id int, request *request.UpdateRoomRequest) error
	DeleteRoom(ctx context.Context, id int) error
}
