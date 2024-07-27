package employeebiz

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/data/response"
	"github.com/qthang02/booking/enities"
	employeerepo "github.com/qthang02/booking/services/employee/repo"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type EmployeeBiz struct {
	employeeRepo employeerepo.IEmployeeRepo
	config       *util.Config
}

func NewEmployeeBiz(employeeRepo employeerepo.IEmployeeRepo, config *util.Config) *EmployeeBiz {
	return &EmployeeBiz{
		employeeRepo: employeeRepo,
		config:       config,
	}
}

func (biz *EmployeeBiz) CreateEmployee(c echo.Context) error {
	var req request.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.CreateEmployee failed to bind create user request")
		return err
	}

	employee, err := biz.employeeRepo.FindEmployeeByEmail(c.Request().Context(), req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msg("EmployeeBiz.validateCreateEmployee Failed to find user by email")
		return err
	}
	if employee != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user with email already exists")
	}

	var user enities.User
	err = copier.Copy(&user, &req)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.CreateEmployee failed to copy copy user")
		return err
	}

	hashPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.CreateEmployee cannot hash password")
		return err
	}

	user.Password = hashPassword

	err = biz.employeeRepo.SaveEmployee(c.Request().Context(), &user)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.CreateEmployee cannot create user")
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (biz *EmployeeBiz) ListEmployees(c echo.Context) error {
	ctx := c.Request().Context()
	paging := &request.Paging{}

	if page, err := strconv.Atoi(c.QueryParam("page")); err == nil {
		paging.Page = page
	}
	if limit, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		paging.Limit = limit
	}

	paging.Process()

	users, err := biz.employeeRepo.ListEmployees(ctx, paging)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.ListEmployees: cannot get users")
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve users")
	}

	var userDTOs []*response.UserDTOResponse
	err = copier.Copy(&userDTOs, &users)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.ListEmployees: cannot copy users")
		return c.JSON(http.StatusInternalServerError, "Failed to copy users")
	}

	resp := response.PaginatedResponse{
		Data:   userDTOs,
		Paging: paging,
	}

	return c.JSON(http.StatusOK, resp)
}

func (biz *EmployeeBiz) GetEmployeeById(c echo.Context) error {
	c.Logger().Info("EmployeeBiz.GetUserById request")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.GetEmployeeById cannot convert id to int")
		return err
	}

	user, err := biz.employeeRepo.FindEmployeeByID(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.GetEmployeeById cannot find user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	userDTO := response.UserDTOResponse{}
	err = copier.Copy(&userDTO, user)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.GetEmployeeById cannot copy user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	return c.JSON(http.StatusOK, userDTO)
}

func (biz *EmployeeBiz) DeleteEmployeeById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.DeleteUserById cannot convert id to int")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	err = biz.employeeRepo.DeleteEmployee(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.DeleteUserById cannot delete user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (biz *EmployeeBiz) UpdateEmployee(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.UpdateEmployee cannot convert id to int")
		return err
	}

	var userUpdateRequest request.UpdateUserRequest

	if err := c.Bind(&userUpdateRequest); err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.UpdateEmployee failed to bind update user request")
		return err
	}

	// validate
	user, err := biz.employeeRepo.FindEmployeeByID(c.Request().Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.UpdateEmployee cannot find user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	if user == nil {
		log.Error().Err(err).Msg("EmployeeBiz.UpdateEmployee cannot find user")
		_ = c.JSON(http.StatusBadRequest, "")
		return err
	}

	err = biz.employeeRepo.UpdateEmployee(c.Request().Context(), id, &userUpdateRequest)
	if err != nil {
		log.Error().Err(err).Msg("EmployeeBiz.UpdateEmployee cannot update user")
		return err
	}

	return c.NoContent(http.StatusOK)
}
