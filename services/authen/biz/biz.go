package biz

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/data/response"
	"github.com/qthang02/booking/enities"
	"github.com/qthang02/booking/services/user/repo"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
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
	var req request.RegisterUserRequest

	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("AuthenBiz.RegisterUser failed to parse request body")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	u, err := biz.userRepo.FindByEmail(c.Request().Context(), req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msg("AuthenBiz.validateRegister Failed to find user by email")
		return err
	}

	if u != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user with email already exists")
	}

	var user enities.User
	if err := copier.Copy(&user, &req); err != nil {
		log.Error().Err(err).Msg("AuthenBiz.RegisterUser failed to copy request")
		_ = c.JSON(http.StatusInternalServerError, "")
		return err
	}

	hashPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("AuthenBiz.RegisterUser cannot hash password")
		return err
	}

	user.Password = hashPassword
	user.Role = util.Customer

	err = biz.userRepo.Save(c.Request().Context(), &user)
	if err != nil {
		log.Error().Err(err).Msg("AuthenBiz.RegisterUser failed to save user")
		_ = c.JSON(http.StatusInternalServerError, "")
		return err
	}

	return c.JSON(http.StatusOK, "")
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (biz *AuthenBiz) Login(c echo.Context) error {
	ctx := c.Request().Context()
	log.Info().Msg("AuthenBiz.Login request")

	var login request.LoginUserRequest
	if err := c.Bind(&login); err != nil {
		log.Error().Err(err).Msg("Failed to bind login request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := biz.userRepo.FindByEmail(ctx, login.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info().Str("email", login.Email).Msg("User not found")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
		log.Error().Err(err).Msg("Failed to find user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	if err := util.VerifyPassword(user.Password, login.Password); err != nil {
		log.Info().Str("email", login.Email).Msg("Invalid password")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	token, err := util.GenerateToken(biz.config.TokenExpiresIn, user, biz.config.TokenSecret)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	resp := LoginResponse{
		Token: token,
	}

	log.Info().Str("email", user.Email).Msg("User logged in successfully")
	return c.JSON(http.StatusOK, resp)
}

func (biz *AuthenBiz) Profile(c echo.Context) error {
	log.Info().Msg("AuthenBiz.Profile request")

	value := c.Get(util.UserID)

	if value == nil {
		log.Info().Msg("AuthenBiz.Profile request has no profile")
	}

	email := fmt.Sprintf("%v", value)

	user, err := biz.userRepo.FindByEmail(c.Request().Context(), email)
	if err != nil {
		log.Error().Err(err).Msg("AuthenBiz.Profile Failed to find user")
		return err
	}

	userDTO := response.UserDTOResponse{}
	err = copier.Copy(&userDTO, user)
	if err != nil {
		log.Error().Err(err).Msg("AuthenBiz.Profile cannot copy user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	return c.JSON(http.StatusOK, userDTO)
}
