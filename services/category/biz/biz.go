package categotybiz

import (
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/data/response"
	categoryrepo "github.com/qthang02/booking/services/category/repo"
	"github.com/qthang02/booking/types"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type CategoryBiz struct {
	repo   categoryrepo.ICategoryRepo
	config *util.Config
}

func NewCategoryBiz(repo categoryrepo.ICategoryRepo, config *util.Config) *CategoryBiz {
	return &CategoryBiz{
		repo:   repo,
		config: config,
	}
}

func (biz *CategoryBiz) List(c echo.Context) error {
	log.Info().Msgf("CategoryBiz.List list categories requset ")

	var paging request.Paging

	err := c.Bind(&paging)
	if err != nil {
		log.Error().Msgf("CategoryBiz.List failed to parse request body err: %v, with requset: %v", err, c.Request())
		return err
	}

	paging.Process()

	categories, err := biz.repo.ListCategories(c.Request().Context(), &paging)
	if err != nil {
		log.Error().Msgf("CategoryBiz.List cannot list catogories error: %v with requset: %v", err, c.Request())
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, category := range categories {
		resp, err := biz.repo.GetCategory(c.Request().Context(), int(category.ID))
		if err != nil {
			log.Error().Msgf("CategoryBiz.List cannot get catogory error: %v with requset: %v", err, c.Request())
			return err
		}
		for _, room := range resp.Rooms {
			if room.Status == types.Ready {
				category.Rooms = append(category.Rooms, room)
			}
		}

		category.AvailableRooms = int64(len(category.Rooms))
	}

	resp := response.ListCategoriesResponse{
		Categories: categories,
		Paging:     &paging,
	}

	return c.JSON(http.StatusOK, resp)
}

func (biz *CategoryBiz) Create(c echo.Context) error {
	log.Info().Msgf("CategoryBiz.Create create category requset ")

	var req request.CreateCategoryRequest

	err := c.Bind(&req)
	if err != nil {
		log.Error().Msgf("CetegoryBiz.Create failed to parse request body err: %v, with requset: %v", err, req)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = biz.repo.CreateCategory(c.Request().Context(), &req)
	if err != nil {
		log.Error().Msgf("CategoryBiz.Create cannot create category error: %v with requset: %v", err, req)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "Create Category Successfully")
}

func (biz *CategoryBiz) Update(c echo.Context) error {
	log.Info().Msgf("CategoryBiz.Update update category requset")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Msgf("CategoryBiz.Update invalid category id error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var req request.UpdateCategoryRequest
	err = c.Bind(&req)
	if err != nil {
		log.Error().Msgf("CategoryBiz.Update failed to parse request body err: %v, with requset: %v", err, req)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	category, err := biz.repo.GetCategory(c.Request().Context(), id)
	if err != nil {
		log.Error().Msgf("CategoryBiz.Update cannot get category error: %v, with id: %d", err, id)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if category == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	err = biz.repo.UpdateCategory(c.Request().Context(), id, &req)
	if err != nil {
		log.Error().Msgf("CategoryBiz.Update cannot update category error: %v with requset: %v", err, req)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Update Category Successfully")
}

func (biz *CategoryBiz) Get(c echo.Context) error {
	log.Info().Msgf("CategoryBiz.GetOrder get category requset ")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Msgf("CategoryBiz.GetOrder invalid category id error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	category, err := biz.repo.GetCategory(c.Request().Context(), id)
	if err != nil {
		log.Error().Msgf("CategoryBiz.GetOrder cannot get category error: %v with id: %d", err, id)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if category == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	category.AvailableRooms = int64(len(category.Rooms))

	return c.JSON(http.StatusOK, category)
}

func (biz *CategoryBiz) Delete(c echo.Context) error {
	log.Info().Msgf("CategoryBiz.Delete delete category requset ")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Msgf("CategoryBiz.Delete invalid category id error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	category, err := biz.repo.GetCategory(c.Request().Context(), id)
	if err != nil {
		log.Error().Msgf("CategoryBiz.Delete cannot get category error: %v with id: %d", err, id)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if category == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	err = biz.repo.DeleteCategory(c.Request().Context(), id)
	if err != nil {
		log.Error().Msgf("CategoryBiz.Delete cannot delete category error: %v with id: %d", err, id)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Delete Category Successfully")
}
