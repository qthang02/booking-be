package categoryrepo

import (
	"github.com/qthang02/booking/enities"
	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) ICategoryRepo {
	db.AutoMigrate(&enities.Category{})

	return &CategoryRepo{
		db: db,
	}
}
