package biz

import (
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/data/response"
	"github.com/qthang02/booking/helper"
	"github.com/qthang02/booking/services/repo"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type UserBiz struct {
	userRepo repo.IUserRepo
	config   *helper.Config
}

func NewUserBiz(userRepo repo.IUserRepo, config *helper.Config) *UserBiz {
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

	hashPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser cannot hash password")
		return err
	}

	user.Password = hashPassword

	err = biz.userRepo.CreateUser(&user)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.CreateUser cannot create user")
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (biz *UserBiz) ListUsers(w http.ResponseWriter, r *http.Request) {
}

func (biz *UserBiz) GetUserById(c echo.Context) error {
	c.Logger().Info("UserBiz.GetUserById request")

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.GetUserById cannot convert id to int")
		return err
	}

	user, err := biz.userRepo.FindByID(idInt)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.GetUserById cannot find user")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	res := response.UserDTOResponse{}
	err = copier.Copy(&res, user)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.GetUserById cannot copy user")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	return c.JSON(http.StatusOK, res)
}

// TODO: delete user

// TODO: update user

func (biz *UserBiz) Login(c echo.Context) error {
	log.Log().Msg("UserBiz.Login request")
	var login requset.LoginUserRequest

	if err := c.Bind(&login); err != nil {
		log.Error().Err(err).Msg("UserBiz.Login failed to bind login request")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	user, err := biz.userRepo.FindByUsername(login.Username)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.Login cannot find user")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	err = helper.VerifyPassword(user.Password, login.Password)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.Login cannot verify password")
		c.JSON(http.StatusUnauthorized, "Username or password is incorrect")
		return err
	}

	token, err := helper.GenerateToken(biz.config.TokenExpiresIn, user.ID, biz.config.TokenSecret)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.Login cannot generate token")
		return err
	}

	return c.JSON(http.StatusOK, token)
}
