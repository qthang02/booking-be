package orderbiz

import (
	"github.com/labstack/echo/v4"
	orderrepo "github.com/qthang02/booking/services/order/repo"
	"github.com/qthang02/booking/util"
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
func (biz *OrderBiz) Get(c echo.Context) error {
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
