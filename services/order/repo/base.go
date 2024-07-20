package orderrepo

import (
	"context"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
)

type IOrderRepo interface {
	SaveOrder(ctx context.Context, order *enities.Order) error
	FindOrder(ctx context.Context, id int) (*enities.Order, error)
	ListOrders(ctx context.Context, paging *request.Paging) ([]*enities.Order, error)
	DeleteOrder(ctx context.Context, id string) error
	UpdateOrder(ctx context.Context, id string, order *enities.Order) error
	CreateOrder(ctx context.Context, order *enities.Order) error
}
