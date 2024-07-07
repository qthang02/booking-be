package userbiz

import (
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/data/response"
	"github.com/qthang02/booking/enities"
	"github.com/qthang02/booking/services/user/repo"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
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
	var user requset.CreateUserRequest

	if err := c.Bind(&user); err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser failed to bind create user request")
		return err
	}

	hashPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser cannot hash password")
		return err
	}

	user.Password = hashPassword

	err = biz.userRepo.CreateUser(c.Request().Context(), &user)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser cannot create user")
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (biz *UserBiz) ListUsers(c echo.Context) error {
	users, err := biz.userRepo.ListUsers(c.Request().Context())
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.ListUsers cannot get all users")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	res := lo.Map(users, func(item *enities.User, _ int) *response.UserDTOResponse {
		var user response.UserDTOResponse
		err := copier.Copy(&user, &item)
		if err != nil {
			log.Error().Err(err).Msgf("UserBiz.ListUsers cannot copy user err: %s", err)
			return nil
		}
		return &user
	})

	return c.JSON(http.StatusOK, res)
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

	res := response.UserDTOResponse{}
	err = copier.Copy(&res, user)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.GetUserById cannot copy user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	return c.JSON(http.StatusOK, res)
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

	var userUpdateRequest requset.UpdateUserRequest

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

	return c.JSON(http.StatusOK, "")
}
