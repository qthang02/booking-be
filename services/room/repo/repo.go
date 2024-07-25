package roomrepo

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

type RoomRepo struct {
	db *gorm.DB
}

func NewRoomRepo(db *gorm.DB) *RoomRepo {
	repo := &RoomRepo{db: db}

	err := db.AutoMigrate(&enities.Room{})
	if err != nil {
		log.Error().Err(err).Msg("NewRoomRepo: failed to auto migrate room database")
		return nil
	}

	err = repo.initRoomDB()
	if err != nil {
		log.Error().Err(err).Msg("NewRoomRepo: failed to init room database")
		return nil
	}

	return repo
}

func (repo *RoomRepo) initRoomDB() error {
	var count int64
	if err := repo.db.Model(&enities.Room{}).Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("initRoomDB: failed to count rooms")
		return err
	}

	if count == 0 {
		log.Info().Msg("initRoomDB: No rooms found, initializing default data")
		for _, room := range database.InitRoomsDataDefault() {
			if err := repo.Save(context.Background(), room); err != nil {
				log.Error().Err(err).Msg("initRoomDB: failed to save room")
				return err
			}
		}
		log.Info().Msg("initRoomDB: Default data initialized successfully")
	} else {
		log.Info().Msgf("initRoomDB: Found %d existing rooms, skipping initialization", count)
	}

	return nil
}

func (repo *RoomRepo) Save(ctx context.Context, room *enities.Room) error {
	err := repo.db.WithContext(ctx).Save(room).Error
	if err != nil {
		log.Error().Err(err).Msg("RoomRepo.Save: failed to save room")
	}
	return err
}

func (repo *RoomRepo) ListRooms(ctx context.Context, paging *request.Paging) ([]*enities.Room, error) {
	var rooms []*enities.Room

	query := repo.db.WithContext(ctx).Model(&enities.Room{})

	if err := query.Count(&paging.Total).Error; err != nil {
		log.Error().Err(err).Msg("RoomRepo.ListRooms: failed to count total rooms")
		return nil, err
	}

	if paging.Limit > 0 {
		query = query.Offset((paging.Page - 1) * paging.Limit).Limit(paging.Limit)
	}

	if err := query.Order("id desc").Find(&rooms).Error; err != nil {
		log.Error().Err(err).Msg("RoomRepo.ListRooms: failed to query rooms")
		return nil, err
	}

	return rooms, nil
}

func (repo *RoomRepo) GetRoom(ctx context.Context, id int) (*enities.Room, error) {
	var room enities.Room
	err := repo.db.WithContext(ctx).First(&room, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Msgf("RoomRepo.GetRoom: room not found with id: %d", id)
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msgf("RoomRepo.GetRoom: failed to get room with id: %d", id)
		return nil, err
	}
	return &room, nil
}

func (repo *RoomRepo) CreateRoom(ctx context.Context, request *request.CreateRoomRequest) error {
	room := enities.Room{}
	err := copier.Copy(&room, request)
	if err != nil {
		log.Error().Err(err).Msg("RoomRepo.CreateRoom: failed to copy request")
		return err
	}

	err = repo.db.WithContext(ctx).Create(&room).Error
	if err != nil {
		log.Error().Err(err).Msg("RoomRepo.CreateRoom: failed to create room")
		return err
	}

	return nil
}

func (repo *RoomRepo) UpdateRoom(ctx context.Context, id int, request *request.UpdateRoomRequest) error {
	var room enities.Room
	err := copier.Copy(&room, request)
	if err != nil {
		log.Error().Err(err).Msg("RoomRepo.UpdateRoom: failed to copy request")
		return err
	}

	result := repo.db.WithContext(ctx).Model(&enities.Room{}).Where("id = ?", id).Updates(&room)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("RoomRepo.UpdateRoom: failed to update room")
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Info().Msgf("RoomRepo.UpdateRoom: room not found with id: %d", id)
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *RoomRepo) DeleteRoom(ctx context.Context, id int) error {
	result := repo.db.WithContext(ctx).Delete(&enities.Room{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("RoomRepo.DeleteRoom: failed to delete room")
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Info().Msgf("RoomRepo.DeleteRoom: room not found with id: %d", id)
		return gorm.ErrRecordNotFound
	}
	return nil
}
