package categoryrepo

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/database"
	"github.com/qthang02/booking/enities"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	categoryRepo *CategoryRepo
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) ICategoryRepo {
	categoryRepo = &CategoryRepo{}

	err := db.AutoMigrate(&enities.Category{})
	if err != nil {
		log.Error().Msgf("NewCategoryRepo: Error migrating category repo: %v", err)
		return nil
	}

	categoryRepo.db = db

	err = categoryRepo.initCategoryDB()
	if err != nil {
		log.Error().Msgf("NewCategoryRepo: Error migrating category repo: %v", err)
		return nil
	}

	return categoryRepo
}

func (repo *CategoryRepo) initCategoryDB() error {
	categories, err := repo.ListCategories(context.Background(), &request.Paging{
		Limit: 10,
		Page:  1,
	})
	if err != nil {
		log.Error().Msgf("initCategoryDB: Error listing categories: %v", err)
		return err
	}

	if len(categories) == 0 {
		for _, category := range database.InitCategoriesDataDefault() {
			err := repo.Save(context.Background(), category)
			if err != nil {
				log.Error().Msgf("initCategoryDB: Error saving category: %v", err)
				return err
			}
		}
	}

	return nil
}

func (repo *CategoryRepo) Save(_ context.Context, category *enities.Category) error {
	err := repo.db.Save(category).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.Save: Error saving category: %v", err)
		return err
	}

	return nil
}

func (repo *CategoryRepo) ListCategories(_ context.Context, paging *request.Paging) ([]*enities.Category, error) {
	var categories []*enities.Category

	err := repo.db.Find(&categories).Count(&paging.Total).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.ListCategories cannot counting result err: %v", err)
		return nil, err
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
	var rooms []enities.Room

	query := "SELECT * FROM categories WHERE id = ?"
	err := repo.db.Raw(query, id).Scan(&category).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.GetCategory cannot find category err: %v with id: %v", err, id)
		return nil, err
	}

	query = "SELECT * FROM rooms WHERE category_id = ?"
	err = repo.db.Raw(query, id).Scan(&rooms).Error
	if err != nil {
		log.Error().Msgf("CategoryRepo.GetCategory cannot preload rooms err: %v with category_id: %v", err, id)
		return nil, err
	}

	category.Rooms = rooms

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
