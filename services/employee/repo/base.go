package employeerepo

import (
	"context"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
)

type IEmployeeRepo interface {
	SaveEmployee(ctx context.Context, employee *enities.User) error
	FindEmployeeByEmail(ctx context.Context, email string) (*enities.User, error)
	FindEmployeeByID(ctx context.Context, id int) (*enities.User, error)
	UpdateEmployee(ctx context.Context, id int, req *request.UpdateUserRequest) error
	DeleteEmployee(ctx context.Context, id int) error
	ListEmployees(ctx context.Context, paging *request.Paging) ([]*enities.User, error)
}
