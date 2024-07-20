package roombiz

import (
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/data/response"
	roomrepo "github.com/qthang02/booking/services/room/repo"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type RoomBiz struct {
	repo   roomrepo.IRoomRepo
	config *util.Config
}

func NewRoomBiz(repo roomrepo.IRoomRepo, config *util.Config) *RoomBiz {
	return &RoomBiz{
		repo:   repo,
		config: config,
	}
}

func (biz *RoomBiz) Get(c echo.Context) error {
	log.Info().Msgf("RoomBiz.GetOrder get room request")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Msgf("RoomBiz.GetOrder room id error: %v, with id: %d", err, id)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	room, err := biz.repo.GetRoom(c.Request().Context(), id)
	if err != nil {
		log.Error().Msgf("RoomBiz.GetOrder room error: %v, with id: %d", err, id)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if room == nil {
		return echo.NewHTTPError(http.StatusNotFound, "room not found")
	}

	return c.JSON(http.StatusOK, room)
}

func (biz *RoomBiz) List(c echo.Context) error {
	log.Info().Msgf("RoomBiz.List get room request")

	var paging request.Paging
	if err := c.Bind(&paging); err != nil {
		log.Error().Msgf("RoomBiz.List room request error: %v, with requset: %v", err, paging)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	paging.Process()

	rooms, err := biz.repo.ListRooms(c.Request().Context(), &paging)
	if err != nil {
		log.Error().Msgf("RoomBiz.List rooms error: %v, with request: %v", err, paging)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := response.ListRoomsResponse{
		Rooms:  rooms,
		Paging: &paging,
	}

	return c.JSON(http.StatusOK, resp)
}

func (biz *RoomBiz) Update(c echo.Context) error {
	log.Info().Msgf("RoomBiz.Update update room request")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Msgf("RoomBiz.Update room id error: %v, with id: %d", err, id)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var req request.UpdateRoomRequest
	if err := c.Bind(&req); err != nil {
		log.Error().Msgf("RoomBiz.Update room request error: %v, with request: %v", err, req)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	room, err := biz.repo.GetRoom(c.Request().Context(), id)
	if err != nil {
		log.Error().Msgf("RoomBiz.Update room error: %v, with id: %d", err, id)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if room == nil {
		return echo.NewHTTPError(http.StatusNotFound, "room not found")
	}

	err = biz.repo.UpdateRoom(c.Request().Context(), id, &req)
	if err != nil {
		log.Error().Msgf("RoomBiz.Update room error: %v, with request: %v", err, req)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Update successfully")
}

func (biz *RoomBiz) Delete(c echo.Context) error {
	log.Info().Msgf("RoomBiz.Delete delete room request")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Msgf("RoomBiz.Delete room id error: %v, with id: %d", err, id)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	room, err := biz.repo.GetRoom(c.Request().Context(), id)
	if err != nil {
		log.Error().Msgf("RoomBiz.Delete room error: %v, with id: %d", err, id)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if room == nil {
		return echo.NewHTTPError(http.StatusNotFound, "room not found")
	}

	err = biz.repo.DeleteRoom(c.Request().Context(), id)
	if err != nil {
		log.Error().Msgf("RoomBiz.Delete room error: %v, with id: %d", err, id)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Delete successfully")
}

func (biz *RoomBiz) Create(c echo.Context) error {
	log.Info().Msgf("RoomBiz.Create create room request")

	var req request.CreateRoomRequest

	if err := c.Bind(&req); err != nil {
		log.Error().Msgf("RoomBiz.Create room request error: %v, with request: %v", err, req)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := biz.repo.CreateRoom(c.Request().Context(), &req)
	if err != nil {
		log.Error().Msgf("RoomBiz.Create room error: %v, with request: %v", err, req)
		return err
	}

	return c.JSON(http.StatusOK, "Create successfully")
}
