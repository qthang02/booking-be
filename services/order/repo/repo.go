package orderrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/database"
	"github.com/qthang02/booking/enities"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderRepo struct {
	db     *gorm.DB
	config *util.Config
}

func NewOrderRepo(db *gorm.DB, config *util.Config) (*OrderRepo, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	repo := &OrderRepo{
		db:     db,
		config: config,
	}

	if err := db.AutoMigrate(&enities.Order{}); err != nil {
		return nil, fmt.Errorf("failed to migrate orders table: %w", err)
	}

	if err := repo.initOrderDB(); err != nil {
		return nil, fmt.Errorf("failed to init orders table: %w", err)
	}

	return repo, nil
}

func (repo *OrderRepo) initOrderDB() error {
	var count int64
	if err := repo.db.Model(&enities.Order{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count orders: %w", err)
	}

	log.Info().Msgf("initOrderDB: found %d existing orders", count)

	if count == 0 {
		for _, order := range database.InitOrdersDataDefault() {
			if err := repo.SaveOrder(context.Background(), order); err != nil {
				return fmt.Errorf("failed to save order: %w", err)
			}
		}
		log.Info().Msg("initOrderDB: default orders created")
	}

	return nil
}

func (repo *OrderRepo) SaveOrder(ctx context.Context, order *enities.Order) error {
	return repo.db.WithContext(ctx).Create(order).Error
}

func (repo *OrderRepo) FindOrder(ctx context.Context, id int) (*enities.Order, error) {
	var order enities.Order
	err := repo.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	return &order, nil
}

func (repo *OrderRepo) ListOrders(ctx context.Context, paging *request.Paging) ([]*enities.Order, error) {
	var orders []*enities.Order

	query := repo.db.WithContext(ctx).Model(&enities.Order{})

	if err := query.Count(&paging.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count orders: %w", err)
	}

	err := query.Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Order("id desc").
		Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, nil
}

func (repo *OrderRepo) DeleteOrder(ctx context.Context, id string) error {
	result := repo.db.WithContext(ctx).Delete(&enities.Order{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete order: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (repo *OrderRepo) UpdateOrder(ctx context.Context, id string, order *enities.Order) error {
	result := repo.db.WithContext(ctx).Model(&enities.Order{}).Where("id = ?", id).Updates(order)
	if result.Error != nil {
		return fmt.Errorf("failed to update order: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (repo *OrderRepo) CreateOrder(ctx context.Context, order *enities.Order) error {
	return repo.db.WithContext(ctx).Create(order).Error
}
