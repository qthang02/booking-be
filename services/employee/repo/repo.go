package employeerepo

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/enities"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type EmployeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) IEmployeeRepo {
	return &EmployeeRepo{db: db}
}

func (repo *EmployeeRepo) SaveEmployee(ctx context.Context, employee *enities.User) error {
	if err := repo.db.WithContext(ctx).Create(employee).Error; err != nil {
		log.Error().Err(err).Msg("EmployeeRepo.SaveEmployee: cannot save employee")
		return err
	}
	return nil
}

func (repo *EmployeeRepo) FindEmployeeByEmail(ctx context.Context, email string) (*enities.User, error) {
	var employee enities.User

	err := repo.db.WithContext(ctx).
		Preload("Orders").
		Where("email = ? AND role IN ?", email, []string{"Admin", "Staff"}).
		First(&employee).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Msgf("EmployeeRepo.FindEmployeeByEmail: Employee not found with email: %v", email)
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msgf("EmployeeRepo.FindEmployeeByEmail: Error fetching employee with email: %v", email)
		return nil, err
	}

	return &employee, nil
}

func (repo *EmployeeRepo) FindEmployeeByID(ctx context.Context, id int) (*enities.User, error) {
	var employee enities.User

	err := repo.db.WithContext(ctx).
		Preload("Orders").
		Where("id = ? AND role IN ?", id, []string{"Admin", "Staff"}).
		First(&employee).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Msgf("EmployeeRepo.FindEmployeeByID: Employee not found with id: %v", id)
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msgf("EmployeeRepo.FindEmployeeByID: Error fetching employee with id: %v", id)
		return nil, err
	}

	return &employee, nil
}

func (repo *EmployeeRepo) UpdateEmployee(ctx context.Context, id int, req *request.UpdateUserRequest) error {
	var employee enities.User

	if err := copier.Copy(&employee, req); err != nil {
		log.Error().Err(err).Msgf("EmployeeRepo.UpdateEmployee: cannot copy employee request")
		return err
	}

	result := repo.db.WithContext(ctx).Model(&enities.User{}).Where("id = ? AND role IN ?", id, []string{"Admin", "Staff"}).Updates(&employee)
	if result.Error != nil {
		log.Error().Err(result.Error).Msgf("EmployeeRepo.UpdateEmployee: cannot update employee")
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Info().Msgf("EmployeeRepo.UpdateEmployee: Employee not found with id: %v", id)
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo *EmployeeRepo) DeleteEmployee(ctx context.Context, id int) error {
	result := repo.db.WithContext(ctx).Where("role IN ?", []string{"Admin", "Staff"}).Delete(&enities.User{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msgf("EmployeeRepo.DeleteEmployee: cannot delete employee")
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Info().Msgf("EmployeeRepo.DeleteEmployee: Employee not found with id: %v", id)
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo *EmployeeRepo) ListEmployees(ctx context.Context, paging *request.Paging) ([]*enities.User, error) {
	var employees []*enities.User

	offset := (paging.Page - 1) * paging.Limit

	var totalCount int64
	if err := repo.db.Model(&enities.User{}).Where("role IN ?", []string{"Admin", "Staff"}).Count(&totalCount).Error; err != nil {
		log.Error().Err(err).Msg("EmployeeRepo.ListEmployees: failed to count total employees")
		return nil, err
	}
	paging.Total = totalCount

	result := repo.db.WithContext(ctx).
		Preload("Orders").
		Where("role IN ?", []string{"Admin", "Staff"}).
		Limit(paging.Limit).
		Offset(offset).
		Find(&employees)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("EmployeeRepo.ListEmployees: failed to list employees")
		return nil, result.Error
	}

	return employees, nil
}
