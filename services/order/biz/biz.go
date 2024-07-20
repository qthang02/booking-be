package orderbiz

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/qthang02/booking/data/request"
	"github.com/qthang02/booking/data/response"
	"github.com/qthang02/booking/enities"
	orderrepo "github.com/qthang02/booking/services/order/repo"
	"github.com/qthang02/booking/util"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type OrderBiz struct {
	orderRepo orderrepo.IOrderRepo
	config    *util.Config
}

func NewOrderBiz(orderRepo orderrepo.IOrderRepo, config *util.Config) *OrderBiz {
	return &OrderBiz{orderRepo: orderRepo, config: config}
}

func (biz *OrderBiz) GetOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	order, err := biz.orderRepo.FindOrder(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, order)
}

func (biz *OrderBiz) CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	var req *request.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if err := validateCreateOrderRequest(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var order enities.Order
	err := copier.Copy(&order, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Call repository to create the order
	err = biz.orderRepo.CreateOrder(ctx, &order)
	if err != nil {
		// Log the error
		log.Error().Err(err).Msg("Failed to create order")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create order"})
	}

	return c.JSON(http.StatusCreated, order)
}

func validateCreateOrderRequest(req *request.CreateOrderRequest) error {
	if req.GuestNumber == 0 {
		return errors.New("guest number must be greater than 0")
	}
	if req.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if req.Checkin.After(req.Checkout) {
		return errors.New("checkin date must be before checkout date")
	}
	return nil
}

func (biz *OrderBiz) UpdateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid order ID"})
	}

	var req request.UpdateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	existingOrder, err := biz.orderRepo.FindOrder(ctx, int(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Order not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve order"})
	}

	if err := validateUpdateOrderRequest(existingOrder); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	err = biz.orderRepo.UpdateOrder(ctx, strconv.FormatUint(id, 10), existingOrder)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update order")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update order"})
	}

	return c.JSON(http.StatusOK, existingOrder)
}

func validateUpdateOrderRequest(order *enities.Order) error {
	if order.GuestNumber == 0 {
		return errors.New("guest number must be greater than 0")
	}
	if order.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if order.Checkin.After(order.Checkout) {
		return errors.New("checkin date must be before checkout date")
	}
	return nil
}

func (biz *OrderBiz) DeleteOrder(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid order ID"})
	}

	_, err = biz.orderRepo.FindOrder(ctx, int(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Order not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve order"})
	}

	err = biz.orderRepo.DeleteOrder(ctx, strconv.FormatUint(id, 10))
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete order")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete order"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Order successfully deleted"})
}

func (biz *OrderBiz) ListOrders(c echo.Context) error {
	ctx := c.Request().Context()

	paging := &request.Paging{
		Page:  1,
		Limit: 10,
	}

	if page, err := strconv.Atoi(c.QueryParam("page")); err == nil {
		paging.Page = page
	}
	if limit, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		paging.Limit = limit
	}

	orders, err := biz.orderRepo.ListOrders(ctx, paging)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list orders")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve orders"})
	}

	resp := response.ListOrdersResponse{
		Orders: orders,
		Paging: paging,
	}

	return c.JSON(http.StatusOK, resp)
}
