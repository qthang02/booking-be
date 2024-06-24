package biz

import (
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/data/response"
	"github.com/qthang02/booking/enities"
	"github.com/qthang02/booking/helper"
	"github.com/qthang02/booking/services/repo"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
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

func (biz *UserBiz) RegisterUser(c echo.Context) error {
	var req requset.RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("UserBiz.RegisterUser failed to parse request body")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	var user enities.User
	if err := copier.Copy(&user, &req); err != nil {
		log.Error().Err(err).Msg("UserBiz.RegisterUser failed to copy request")
		c.JSON(http.StatusInternalServerError, "")
		return err
	}

	hashPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.RegisterUser cannot hash password")
		return err
	}

	user.Password = hashPassword

	err = biz.userRepo.Save(&user)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.RegisterUser failed to save user")
		c.JSON(http.StatusInternalServerError, "")
		return err
	}

	return c.JSON(http.StatusOK, "")
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

func (biz *UserBiz) ListUsers(c echo.Context) error {
	users, err := biz.userRepo.GetAllUsers()
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.ListUsers cannot get all users")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	res := lo.Map(users, func(item *enities.User, _ int) *response.UserDTOResponse {
		var user response.UserDTOResponse
		copier.Copy(&user, &item)
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

	user, err := biz.userRepo.FindByID(id)
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

func (biz *UserBiz) DeleteUserById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.DeleteUserById cannot convert id to int")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	err = biz.userRepo.DeleteUser(id)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.DeleteUserById cannot delete user")
		c.JSON(http.StatusBadRequest, "")
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
	user, err := biz.userRepo.FindByID(id)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.UpdateUser cannot find user")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	if user == nil {
		log.Error().Err(err).Msg("UserBiz.UpdateUser cannot find user")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	err = biz.userRepo.UpdateUser(id, &userUpdateRequest)
	if err != nil {
		log.Error().Err(err).Msg("UserBiz.UpdateUser cannot update user")
		return err
	}

	return c.JSON(http.StatusOK, "")
}

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
