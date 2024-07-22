package userbiz

import (
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/data/response"
	"github.com/qthang02/booking/enities"
	"github.com/qthang02/booking/services/user/repo"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type UserBiz struct {
	userRepo userrepo.IUserRepo
	config   *util.Config
}

func NewUserBiz(userRepo userrepo.IUserRepo, config *util.Config) *UserBiz {
	return &UserBiz{
		userRepo: userRepo,
		config:   config,
	}
}

func (biz *UserBiz) CreateUser(c echo.Context) error {
	var req request.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser failed to bind create user request")
		return err
	}

	var user enities.User
	err := copier.Copy(&user, &req)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser failed to copy copy user")
		return err
	}

	hashPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser cannot hash password")
		return err
	}

	user.Password = hashPassword

	err = biz.userRepo.Save(c.Request().Context(), &user)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser cannot create user")
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (biz *UserBiz) ListUsers(c echo.Context) error {
	ctx := c.Request().Context()
	paging := &request.Paging{}

	if page, err := strconv.Atoi(c.QueryParam("page")); err == nil {
		paging.Page = page
	}
	if limit, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		paging.Limit = limit
	}

	paging.Process()

	users, err := biz.userRepo.ListUsers(ctx, paging)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.ListUsers: cannot get users")
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve users")
	}

	var userDTOs []*response.UserDTOResponse
	err = copier.Copy(&userDTOs, &users)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.ListUsers: cannot copy users")
		return c.JSON(http.StatusInternalServerError, "Failed to copy users")
	}

	resp := response.PaginatedResponse{
		Data:   userDTOs,
		Paging: paging,
	}

	return c.JSON(http.StatusOK, resp)
}

func (biz *UserBiz) GetUserById(c echo.Context) error {
	c.Logger().Info("UserBiz.GetUserById request")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.GetUserById cannot convert id to int")
		return err
	}

	user, err := biz.userRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.GetUserById cannot find user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	userDTO := response.UserDTOResponse{}
	err = copier.Copy(&userDTO, user)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.GetUserById cannot copy user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	resp := response.ProfileResponse{
		User: userDTO,
	}

	return c.JSON(http.StatusOK, resp)
}

func (biz *UserBiz) DeleteUserById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.DeleteUserById cannot convert id to int")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	err = biz.userRepo.DeleteUser(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.DeleteUserById cannot delete user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (biz *UserBiz) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.UpdateUser cannot convert id to int")
		return err
	}

	var userUpdateRequest request.UpdateUserRequest

	if err := c.Bind(&userUpdateRequest); err != nil {
		log.Error().Err(err).Msg("UserBiz.UpdateUser failed to bind update user request")
		return err
	}

	// validate
	user, err := biz.userRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.UpdateUser cannot find user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	if user == nil {
		log.Error().Err(err).Msg("UserBiz.UpdateUser cannot find user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	err = biz.userRepo.UpdateUser(c.Request().Context(), id, &userUpdateRequest)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.UpdateUser cannot update user")
		return err
	}

	return c.NoContent(http.StatusOK)
}
