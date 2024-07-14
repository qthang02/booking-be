package roomrepo

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
	roomRepo *RoomRepo
)

type RoomRepo struct {
	db *gorm.DB
}

func NewRoomRepo(db *gorm.DB) *RoomRepo {
	roomRepo = &RoomRepo{}

	err := db.AutoMigrate(&enities.Room{})
	if err != nil {
		log.Fatal().Err(err).Msg("NewRoomRepo: failed to auto migrate room database")
		return nil
	}

	roomRepo.db = db

	err = roomRepo.initRoomDB()
	if err != nil {
		log.Fatal().Err(err).Msg("NewRoomRepo: failed to init room database")
		return nil
	}

	return roomRepo
}

func (repo *RoomRepo) initRoomDB() error {
	rooms, err := repo.ListRooms(context.Background(), &request.Paging{})
	if err != nil {
		log.Fatal().Err(err).Msg("initRoomDB: failed to list rooms")
		return err
	}

	if len(rooms) == 0 {
		for _, room := range database.InitRoomsDataDefault() {
			err := repo.Save(context.Background(), room)
			if err != nil {
				log.Fatal().Err(err).Msg("initRoomDB: failed to save room")
				return err
			}
		}
	}

	return nil
}

func (repo *RoomRepo) Save(_ context.Context, room *enities.Room) error {
	err := repo.db.Save(room).Error
	if err != nil {
		log.Error().Err(err).Msg("RoomRepo.Save: failed to save room")
	}

	return err
}

func (repo *RoomRepo) ListRooms(_ context.Context, paging *request.Paging) ([]*enities.Room, error) {
	var rooms []*enities.Room

	err := repo.db.Find(&rooms).Count(&paging.Total).Error
	if err != nil {
		log.Error().Msgf("RoomRepo.ListRooms count total error: %v", err)
		return nil, err
	}

	repo.db = repo.db.Offset((paging.Page - 1) * paging.Limit)

	err = repo.db.Order("id desc").Find(&rooms).Limit(paging.Limit).Error
	if err != nil {
		log.Error().Msgf("RoomRepo.ListRooms limit error: %v", err)
		return nil, err
	}

	return rooms, nil
}

func (repo *RoomRepo) GetRoom(_ context.Context, id int) (*enities.Room, error) {
	log.Info().Msgf("RoomRepo.GetRoom get room request: %v", id)

	var room enities.Room
	err := repo.db.Raw("select * from rooms where id = ?", id).Scan(&room).Error
	if err != nil {
		log.Error().Msgf("RoomRepo.GetRoom find room error: %v", err)
		return nil, err
	}

	return &room, nil
}

func (repo *RoomRepo) CreateRoom(_ context.Context, request *request.CreateRoomRequest) error {
	log.Info().Msgf("RoomRepo.CreateRoom create room request: %v", request)

	room := enities.Room{}
	err := copier.Copy(&room, request)
	if err != nil {
		log.Error().Msgf("RoomRepo.CreateRoom cannot copy request error: %v", err)
		return err
	}

	err = repo.db.Create(&room).Error
	if err != nil {
		log.Error().Msgf("RoomRepo.CreateRoom create room error: %v", err)
		return err
	}

	return nil
}

func (repo *RoomRepo) UpdateRoom(_ context.Context, id int, request *request.UpdateRoomRequest) error {
	log.Info().Msgf("RoomRepo.UpdateRoom update room request: %v", request)

	var room enities.Room
	err := copier.Copy(&room, request)
	if err != nil {
		log.Error().Msgf("RoomRepo.UpdateRoom cannot copy request error: %v", err)
		return err
	}

	err = repo.db.Where("id = ?", id).Updates(&room).Error
	if err != nil {
		log.Error().Msgf("RoomRepo.UpdateRoom update room error: %v", err)
		return err
	}
	return nil
}

func (repo *RoomRepo) DeleteRoom(_ context.Context, id int) error {
	log.Info().Msgf("RoomRepo.DeleteRoom delete room request: %v", id)

	err := repo.db.Delete(&enities.Room{}, id).Error
	if err != nil {
		log.Error().Msgf("RoomRepo.DeleteRoom delete room error: %v", err)
		return err
	}

	return nil
}
