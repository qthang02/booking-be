package categoryrepo

import (
	"context"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
)

type ICategoryRepo interface {
	Save(ctx context.Context, category *enities.Category) error
	ListCategories(ctx context.Context, paging *request.Paging) ([]*enities.Category, error)
	GetCategory(ctx context.Context, id int) (*enities.Category, error)
	DeleteCategory(ctx context.Context, id int) error
	CreateCategory(ctx context.Context, request *request.CreateCategoryRequest) error
	UpdateCategory(ctx context.Context, id int, request *request.UpdateCategoryRequest) error
}
