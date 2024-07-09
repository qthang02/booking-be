package categoryrepo

import (
	"context"
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
)

type ICategoryRepo interface {
	ListCategories(ctx context.Context, paging *requset.Paging) ([]*enities.Category, error)
	GetCategory(ctx context.Context, id int) (*enities.Category, error)
	DeleteCategory(ctx context.Context, id int) error
	CreateCategory(ctx context.Context, request *requset.CreateCategoryRequest) error
	UpdateCategory(ctx context.Context, id int, request *requset.UpdateCategoryRequest) error
}
