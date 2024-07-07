package categoryrepo

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) ICategoryRepo {
	err := db.AutoMigrate(&enities.Category{})
	if err != nil {
		log.Error().Msgf("Error migrating category repo: %v", err)
		return nil
	}

	return &CategoryRepo{
		db: db,
	}
}

func (repo *CategoryRepo) ListCategories(_ context.Context, paging *requset.Paging) ([]*enities.Category, error) {
	log.Info().Msgf("CategoryRepo.ListCategories Listing categories with paging: %v", paging)

	var categories []*enities.Category

	err := repo.db.Find(&categories).Count(&paging.Total).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.ListCategories cannot counting result err: %v", err)
	}

	repo.db = repo.db.Offset((paging.Page - 1) * paging.Limit)

	err = repo.db.Order("id desc").Find(&categories).Limit(paging.Limit).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.ListCategories cannot limiting result err: %v", err)
		return nil, err
	}

	return categories, nil
}

func (repo *CategoryRepo) GetCategory(_ context.Context, id int) (*enities.Category, error) {
	log.Info().Msgf("CategoryRepo.GetCategory Getting category with id: %v", id)

	var category enities.Category

	err := repo.db.First(&category, id).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.GetCategory cannot find category with id: %v", id)
		return nil, err
	}

	return &category, nil
}

func (repo *CategoryRepo) DeleteCategory(_ context.Context, id int) error {
	log.Info().Msgf("CategoryRepo.DeleteCategory Delete category with id: %v", id)
	err := repo.db.Delete(&enities.Category{}, "id = ?", id).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.DeleteCategory cannot delete category with id: %v", id)
		return err
	}

	return nil
}

func (repo *CategoryRepo) CreateCategory(_ context.Context, request *requset.CreateCategoryRequest) error {
	log.Info().Msgf("CategoryRepo.CreateCategory Create category request: %v", request)

	category := enities.Category{}
	err := copier.Copy(&category, request)
	if err != nil {
		log.Error().Msgf("CategoryRepo.CreateCategory cannot copy request: %v", err)
		return err
	}

	err = repo.db.Create(&category).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.CreateCategory cannot create category with id: %v", category.ID)
		return err
	}

	return nil
}

func (repo *CategoryRepo) UpdateCategory(_ context.Context, request requset.UpdateCategoryRequest) error {
	log.Info().Msgf("CategoryRepo.UpdateCategory Update category request: %v", request)
	category := enities.Category{}
	err := copier.Copy(&category, request)
	if err != nil {
		log.Error().Msgf("CategoryRepo.UpdateCategory cannot copy request: %v", err)
		return err
	}

	err = repo.db.Save(&category).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.UpdateCategory cannot update category with id: %v", category.ID)
		return err
	}

	return nil
}
