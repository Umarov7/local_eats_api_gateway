package handler

import (
	pb "api-gateway/genproto/order"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateOrder godoc
// @Summary Creates an order
// @Description Inserts a new order into database
// @Tags order
// @Security ApiKeyAuth
// @Param order body order.NewOrder true "Order info"
// @Success 200 {object} order.NewOrderResp
// @Failure 400 {object} string "Invalid order data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	h.Logger.Info("CreateOrder method is starting")

	var data pb.NewOrder
	if err := c.ShouldBindJSON(&data); err != nil {
		er := errors.Wrap(err, "invalid order data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.OrderClient.MakeOrder(ctx, &data)
	if err != nil {
		er := errors.Wrap(err, "error creating order").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("Order created successfully")
	c.JSON(http.StatusOK, res)
}

// GetOrderByID godoc
// @Summary Gets an order
// @Description Gets order from database
// @Tags order
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Success 200 {object} order.OrderInfo
// @Failure 400 {object} string "Invalid order ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /orders/{id} [get]
func (h *Handler) GetOrderByID(c *gin.Context) {
	h.Logger.Info("GetOrderByID method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid order id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.OrderClient.GetOrderByID(ctx, &pb.ID{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting order").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("GetOrderByID method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// ChangeStatus godoc
// @Summary Updates an order
// @Description Updates order status in database
// @Tags order
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param status body order.StatusNoID true "Order status"
// @Success 200 {object} order.UpdatedOrder
// @Failure 400 {object} string "Invalid order ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /orders/{id}/status [put]
func (h *Handler) ChangeStatus(c *gin.Context) {
	h.Logger.Info("ChangeStatus method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid order id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	var data pb.StatusNoID
	if err := c.ShouldBindJSON(&data); err != nil {
		er := errors.Wrap(err, "invalid order data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.OrderClient.ChangeStatus(ctx, &pb.Status{
		Id:     id,
		Status: data.Status,
	})
	if err != nil {
		er := errors.Wrap(err, "error changing order status").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("ChangeStatus method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// FetchOrdersForCustomer godoc
// @Summary Gets orders for customer
// @Description Gets orders from database
// @Tags order
// @Security ApiKeyAuth
// @Param page query int true "Page number"
// @Param limit query int true "Number of items per page"
// @Success 200 {object} order.OrdersCustomer
// @Failure 500 {object} string "Server error while processing request"
// @Router /orders [get]
func (h *Handler) FetchOrdersForCustomer(c *gin.Context) {
	h.Logger.Info("FetchOrdersForCustomer method is starting")

	page := c.Query("page")
	limit := c.Query("limit")

	p, err := strconv.Atoi(page)
	if err != nil {
		er := errors.Wrap(err, "invalid pagination parameters").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	l, err := strconv.Atoi(limit)
	if err != nil {
		er := errors.Wrap(err, "invalid pagination parameters").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.OrderClient.FetchOrdersForCustomer(ctx, &pb.Pagination{
		Limit:  int32(l),
		Offset: int32((p - 1) * l),
	})
	if err != nil {
		er := errors.Wrap(err, "error getting orders").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("FetchOrdersForCustomer method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// FetchOrdersForKitchen godoc
// @Summary Gets orders for kitchen
// @Description Gets orders from database
// @Tags order
// @Security ApiKeyAuth
// @Param id path string true "Kitchen ID"
// @Param status query string true "Status"
// @Param page query int true "Page number"
// @Param limit query int true "Number of items per page"
// @Success 200 {object} order.OrdersKitchen
// @Failure 500 {object} string "Server error while processing request"
// @Router /kitchens/{id}/orders [get]
func (h *Handler) FetchOrdersForKitchen(c *gin.Context) {
	h.Logger.Info("FetchOrdersForKitchen method is starting")

	kitchenID := c.Param("id")
	status := c.Query("status")
	page := c.Query("page")
	limit := c.Query("limit")

	_, err := uuid.Parse(kitchenID)
	if err != nil {
		er := errors.Wrap(err, "invalid dish ID").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		er := errors.Wrap(err, "invalid pagination parameters").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	l, err := strconv.Atoi(limit)
	if err != nil {
		er := errors.Wrap(err, "invalid pagination parameters").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.OrderClient.FetchOrdersForKitchen(ctx, &pb.Filter{
		KitchenId: kitchenID,
		Status:    status,
		Pagination: &pb.Pagination{
			Limit:  int32(l),
			Offset: int32((p - 1) * l),
		},
	})
	if err != nil {
		er := errors.Wrap(err, "error getting orders").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("FetchOrdersForKitchen method has finished successfully")
	c.JSON(http.StatusOK, res)
}
