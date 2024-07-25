package categoryrepo

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/database"
	"github.com/qthang02/booking/enities"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	repo := &CategoryRepo{db: db}

	err := db.AutoMigrate(&enities.Category{})
	if err != nil {
		log.Error().Err(err).Msg("NewCategoryRepo: Error migrating category repo")
		return nil
	}

	err = repo.initCategoryDB()
	if err != nil {
		log.Error().Err(err).Msg("NewCategoryRepo: Error initializing category repo")
		return nil
	}

	return repo
}

func (repo *CategoryRepo) initCategoryDB() error {
	var count int64
	if err := repo.db.Model(&enities.Category{}).Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("initCategoryDB: Error counting categories")
		return err
	}

	if count == 0 {
		log.Info().Msg("initCategoryDB: No categories found, initializing default data")
		for _, category := range database.InitCategoriesDataDefault() {
			if err := repo.Save(context.Background(), category); err != nil {
				log.Error().Err(err).Msg("initCategoryDB: Error saving category")
				return err
			}
		}
		log.Info().Msg("initCategoryDB: Default categories initialized successfully")
	} else {
		log.Info().Msgf("initCategoryDB: Found %d existing categories, skipping initialization", count)
	}

	return nil
}

func (repo *CategoryRepo) Save(ctx context.Context, category *enities.Category) error {
	err := repo.db.WithContext(ctx).Save(category).Error
	if err != nil {
		log.Error().Err(err).Msg("CategoryRepo.Save: Error saving category")
	}
	return err
}

func (repo *CategoryRepo) ListCategories(ctx context.Context, paging *request.Paging) ([]*enities.Category, error) {
	var categories []*enities.Category

	query := repo.db.WithContext(ctx).Model(&enities.Category{})

	if err := query.Count(&paging.Total).Error; err != nil {
		log.Error().Err(err).Msg("CategoryRepo.ListCategories: Error counting categories")
		return nil, err
	}

	err := query.Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Order("id desc").
		Find(&categories).Error
	if err != nil {
		log.Error().Err(err).Msg("CategoryRepo.ListCategories: Error fetching categories")
		return nil, err
	}

	return categories, nil
}

func (repo *CategoryRepo) GetCategory(ctx context.Context, id int) (*enities.Category, error) {
	var category enities.Category

	err := repo.db.WithContext(ctx).Preload("Rooms").First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Msgf("CategoryRepo.GetCategory: Category not found with id: %d", id)
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msgf("CategoryRepo.GetCategory: Error fetching category with id: %d", id)
		return nil, err
	}

	return &category, nil
}

func (repo *CategoryRepo) DeleteCategory(ctx context.Context, id int) error {
	result := repo.db.WithContext(ctx).Delete(&enities.Category{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msgf("CategoryRepo.DeleteCategory: Error deleting category with id: %d", id)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Info().Msgf("CategoryRepo.DeleteCategory: Category not found with id: %d", id)
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *CategoryRepo) CreateCategory(ctx context.Context, request *request.CreateCategoryRequest) error {
	category := enities.Category{}
	if err := copier.Copy(&category, request); err != nil {
		log.Error().Err(err).Msg("CategoryRepo.CreateCategory: Error copying request")
		return err
	}

	if err := repo.db.WithContext(ctx).Create(&category).Error; err != nil {
		log.Error().Err(err).Msg("CategoryRepo.CreateCategory: Error creating category")
		return err
	}

	return nil
}

func (repo *CategoryRepo) UpdateCategory(ctx context.Context, id int, request *request.UpdateCategoryRequest) error {
	category := enities.Category{}
	if err := copier.Copy(&category, request); err != nil {
		log.Error().Err(err).Msg("CategoryRepo.UpdateCategory: Error copying request")
		return err
	}

	result := repo.db.WithContext(ctx).Model(&enities.Category{}).Where("id = ?", id).Updates(&category)
	if result.Error != nil {
		log.Error().Err(result.Error).Msgf("CategoryRepo.UpdateCategory: Error updating category with id: %d", id)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Info().Msgf("CategoryRepo.UpdateCategory: Category not found with id: %d", id)
		return gorm.ErrRecordNotFound
	}

	return nil
}
