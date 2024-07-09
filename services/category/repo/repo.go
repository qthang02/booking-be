package categoryrepo

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/qthang02/booking/data/request"
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

func (repo *CategoryRepo) ListCategories(_ context.Context, paging *request.Paging) ([]*enities.Category, error) {
	log.Info().Msgf("CategoryRepo.ListCategories Listing categories with paging: %v", paging)

	var categories []*enities.Category

	err := repo.db.Find(&categories).Count(&paging.Total).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.ListCategories cannot counting result err: %v", err)
		return nil, err
	}

	repo.db = repo.db.Offset((paging.Page - 1) * paging.Limit)

	err = repo.db.Preload("Rooms").First(&categories).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.ListCategories cannot find result err: %v", err)
		return nil, err
	}

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

	err := repo.db.Where("id = ?", id).Preload("Rooms").First(&category).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.GetCategory cannot find category err: %v with id: %v", err, id)
		return nil, err
	}

	err = repo.db.Find(&category.Rooms).Count(&category.AvailableRooms).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.GetCategory cannot count rooms in category err: %v with id: %v", err, id)
		return nil, err
	}

	return &category, nil
}

func (repo *CategoryRepo) DeleteCategory(_ context.Context, id int) error {
	log.Info().Msgf("CategoryRepo.DeleteCategory Delete category with id: %v", id)
	err := repo.db.Delete(&enities.Category{}, "id = ?", id).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.DeleteCategory cannot delete category err: %v with id: %v", err, id)
		return err
	}

	return nil
}

func (repo *CategoryRepo) CreateCategory(_ context.Context, request *request.CreateCategoryRequest) error {
	log.Info().Msgf("CategoryRepo.CreateCategory Create category request: %v", request)

	category := enities.Category{}
	err := copier.Copy(&category, request)
	if err != nil {
		log.Error().Msgf("CategoryRepo.CreateCategory cannot copy request: %v", err)
		return err
	}

	err = repo.db.Create(&category).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.CreateCategory cannot create category err: %v with id: %v", err, category.ID)
		return err
	}

	return nil
}

func (repo *CategoryRepo) UpdateCategory(_ context.Context, id int, request *request.UpdateCategoryRequest) error {
	log.Info().Msgf("CategoryRepo.UpdateCategory Update category request: %v", request)
	category := enities.Category{}
	err := copier.Copy(&category, request)
	if err != nil {
		log.Error().Msgf("CategoryRepo.UpdateCategory cannot copy request: %v", err)
		return err
	}

	err = repo.db.Where(id).Updates(&category).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.UpdateCategory cannot update category err: %v with id: %v", err, category.ID)
		return err
	}

	return nil
}
