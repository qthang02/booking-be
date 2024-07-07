package biz

import (
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/requset"
	"github.com/qthang02/booking/enities"
	"github.com/qthang02/booking/services/user/repo"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"net/http"
)

type AuthenBiz struct {
	userRepo userrepo.IUserRepo
	config   *util.Config
}

func NewAuthenBiz(userRepo userrepo.IUserRepo, config *util.Config) *AuthenBiz {
	return &AuthenBiz{
		userRepo: userRepo,
		config:   config,
	}
}

func (biz *AuthenBiz) RegisterUser(c echo.Context) error {
	var req requset.RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("AuthenBiz.RegisterUser failed to parse request body")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	var user enities.User
	if err := copier.Copy(&user, &req); err != nil {
		log.Error().Err(err).Msg("AuthenBiz.RegisterUser failed to copy request")
		c.JSON(http.StatusInternalServerError, "")
		return err
	}

	hashPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("AuthenBiz.RegisterUser cannot hash password")
		return err
	}

	user.Password = hashPassword

	err = biz.userRepo.Save(&user)
	if err != nil {
		log.Error().Err(err).Msg("AuthenBiz.RegisterUser failed to save user")
		c.JSON(http.StatusInternalServerError, "")
		return err
	}

	return c.JSON(http.StatusOK, "")
}

func (biz *AuthenBiz) Login(c echo.Context) error {
	log.Log().Msg("AuthenBiz.Login request")
	var login requset.LoginUserRequest

	if err := c.Bind(&login); err != nil {
		log.Error().Err(err).Msg("AuthenBiz.Login failed to bind login request")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	user, err := biz.userRepo.FindByEmail(login.Email)
	if err != nil {
		log.Error().Err(err).Msg("AuthenBiz.Login cannot find user")
		c.JSON(http.StatusBadRequest, "")
		return err
	}

	err = util.VerifyPassword(user.Password, login.Password)
	if err != nil {
		log.Error().Err(err).Msg("AuthenBiz.Login cannot verify password")
		c.JSON(http.StatusUnauthorized, "Username or password is incorrect")
		return err
	}

	token, err := util.GenerateToken(biz.config.TokenExpiresIn, user.ID, biz.config.TokenSecret)
	if err != nil {
		log.Error().Err(err).Msg("AuthenBiz.Login cannot generate token")
		return err
	}

	return c.JSON(http.StatusOK, token)
}
