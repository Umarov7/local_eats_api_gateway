package handler

import (
	pb "api-gateway/genproto/payment"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreatePayment godoc
// @Summary Creates a payment
// @Description Inserts a new payment into database
// @Tags payment
// @Security ApiKeyAuth
// @Param payment body payment.NewPayment true "Payment info"
// @Success 200 {object} payment.NewPayment
// @Failure 400 {object} string "Invalid payment data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /payments [post]
func (h *Handler) CreatePayment(c *gin.Context) {
	h.Logger.Info("CreatePayment method is starting")

	var data pb.NewPayment
	if err := c.ShouldBindJSON(&data); err != nil {
		er := errors.Wrap(err, "invalid payment data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	if data.CardNumber != "" {
		if len(data.CardNumber) != 16 {
			er := errors.New("invalid card number").Error()
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": er})
			h.Logger.Error(er)
			return
		}
	}

	if data.ExpiryDate != "" {
		if len(data.ExpiryDate) != 5 {
			er := errors.New("invalid expiry date").Error()
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": er})
			h.Logger.Error(er)
			return
		}
	}

	if data.Cvv != "" {
		if len(data.Cvv) != 3 {
			er := errors.New("invalid CVV").Error()
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": er})
			h.Logger.Error(er)
			return
		}
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.PaymentClient.MakePayment(ctx, &data)
	if err != nil {
		er := errors.Wrap(err, "error creating payment").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetPayment godoc
// @Summary Gets a payment
// @Description Retrieves payment info from database
// @Tags payment
// @Security ApiKeyAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} payment.PaymentDetails
// @Failure 400 {object} string "Invalid payment ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /payments/{id} [get]
func (h *Handler) GetPayment(c *gin.Context) {
	h.Logger.Info("GetPayment method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid payment id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.PaymentClient.GetPayment(ctx, &pb.ID{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting payment").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	c.JSON(http.StatusOK, res)
}
