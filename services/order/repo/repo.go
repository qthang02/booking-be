package orderrepo

import (
	"context"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/database"
	"github.com/qthang02/booking/enities"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	orderRepo *OrderRepo
)

type OrderRepo struct {
	db     *gorm.DB
	config *util.Config
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	orderRepo = &OrderRepo{}

	err := db.AutoMigrate(&enities.Order{})
	if err != nil {
		log.Error().Msgf("NewOrderRepo: failed to migrate orders table: %v", err)
		return nil
	}

	orderRepo.db = db

	err = orderRepo.initOrderDB()
	if err != nil {
		log.Error().Msgf("NewOrderRepo: failed to init orders table: %v", err)
		return nil
	}

	return orderRepo
}

func (repo *OrderRepo) initOrderDB() error {
	orders, err := repo.ListOrders(context.Background(), &request.Paging{})
	if err != nil {
		log.Error().Msgf("OrderRepo: failed to list orders: %v", err)
		return err
	}

	if len(orders) == 0 {
		for _, order := range database.InitOrdersDataDefault() {
			err := repo.SaveOrder(context.Background(), order)
			if err != nil {
				log.Error().Msgf("OrderRepo: failed to save order: %v", err)
				return err
			}
		}
	}

	return nil
}

func (repo *OrderRepo) SaveOrder(_ context.Context, order *enities.Order) error {
	err := repo.db.Create(order).Error
	if err != nil {
		log.Error().Msgf("SaveOrder: failed to save order: %v", err)
	}

	return err
}

func (repo *OrderRepo) FindOrder(_ context.Context, id int) (*enities.Order, error) {
	order := &enities.Order{}
	err := repo.db.Where("id = ?", id).First(order).Error
	if err != nil {
		log.Error().Msgf("FindOrder: failed to find order: %v", err)
	}

	return order, err
}

func (repo *OrderRepo) ListOrders(ctx context.Context, paging *request.Paging) ([]*enities.Order, error) {
	var orders []*enities.Order

	err := repo.db.Find(&orders).Count(&paging.Total).Error
	if err != nil {
		log.Error().Msgf("ListOrders: failed to list orders: %v", err)
		return nil, err
	}

	repo.db = repo.db.Offset((paging.Page - 1) * paging.Limit)

	err = repo.db.Order("id desc").Find(&orders).Limit(paging.Limit).Error
	if err != nil {
		log.Error().Msgf("RoomRepo.ListRooms limit error: %v", err)
		return nil, err
	}

	return orders, nil
}

func (repo *OrderRepo) DeleteOrder(_ context.Context, id string) error {
	order := &enities.Order{}
	err := repo.db.Where("id = ?", id).Delete(order).Error
	if err != nil {
		log.Error().Msgf("DeleteOrder: failed to delete order: %v", err)
	}

	return err
}

func (repo *OrderRepo) UpdateOrder(_ context.Context, id string, order *enities.Order) error {
	err := repo.db.Model(&enities.Order{}).Where("id = ?", id).Updates(order).Error
	if err != nil {
		log.Error().Msgf("UpdateOrder: failed to update order: %v", err)
	}

	return err
}

func (repo *OrderRepo) CreateOrder(_ context.Context, order *enities.Order) error {
	err := repo.db.Create(order).Error
	if err != nil {
		log.Error().Msgf("CreateOrder: failed to save order: %v", err)
	}

	return err
}
