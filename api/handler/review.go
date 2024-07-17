package handler

import (
	pb "api-gateway/genproto/review"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateReview godoc
// @Summary Creates a review
// @Description Inserts a new review into database
// @Tags review
// @Security ApiKeyAuth
// @Param review body review.NewReview true "Review info"
// @Success 200 {object} review.NewReviewResp
// @Failure 400 {object} string "Invalid review data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /reviews [post]
func (h *Handler) CreateReview(c *gin.Context) {
	h.Logger.Info("CreateReview method is starting")

	var data pb.NewReview
	if err := c.ShouldBindJSON(&data); err != nil {
		er := errors.Wrap(err, "invalid review data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.ReviewClient.RateAndComment(ctx, &data)
	if err != nil {
		er := errors.Wrap(err, "failed to create review").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetReviews godoc
// @Summary Gets reviews
// @Description Gets reviews from database
// @Tags review
// @Security ApiKeyAuth
// @Param id path string true "Kitchen ID"
// @Param page query int true "Page number"
// @Param limit query int true "Number of items per page"
// @Success 200 {object} review.Reviews
// @Failure 400 {object} string "Invalid review data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /kitchens/{id}/reviews [get]
func (h *Handler) GetReviews(c *gin.Context) {
	h.Logger.Info("GetReviews method is starting")

	kitchenID := c.Param("id")
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

	res, err := h.ReviewClient.GetReviewOfKitchen(ctx, &pb.Filter{
		KitchenId: kitchenID,
		Limit:     int32(l),
		Offset:    int32((p - 1) * l),
	})
	if err != nil {
		er := errors.Wrap(err, "error getting reviews").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	c.JSON(http.StatusOK, res)
}
