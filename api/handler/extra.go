package handler

import (
	pb "api-gateway/genproto/extra"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// GetStatistics godoc
// @Summary Gets kitchen's statistics
// @Description Informs about kitchen statistics by date
// @Tags kitchen
// @Param kitchen_id path string true "Kitchen ID"
// @Success 200 {object} extra.Statistics
// @Failure 400 {object} string "Invalid kitchen ID or date format"
// @Failure 500 {object} string "Server error while processing request"
// @Router /local-eats/kitchens/{id}/statistics [get]
func (h *Handler) GetStatistics(c *gin.Context) {
	h.Logger.Info("GetStatistics method is starting")
	kitchenID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	_, err := uuid.Parse(kitchenID)
	if err != nil {
		er := errors.Wrap(err, "invalid kitchen id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	_, err = time.Parse("2006-01-02", startDate)
	if err != nil {
		er := errors.Wrap(err, "invalid start date").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		er := errors.Wrap(err, "invalid end date").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := h.ExtraClient.GetStatistics(ctx, &pb.Period{
		Id:        kitchenID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		er := errors.Wrap(err, "error getting statistics").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("GetStatistics method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// TrackActivity godoc
// @Summary Tracks user's activity
// @Description Informs about user's activity by date
// @Tags user
// @Param user_id path string true "User ID"
// @Success 200 {object} extra.Activity
// @Failure 400 {object} string "Invalid user ID or date format"
// @Failure 500 {object} string "Server error while processing request"
// @Router /local-eats/users/{id}/activity [get]
func (h *Handler) TrackActivity(c *gin.Context) {
	h.Logger.Info("TrackActivity method is starting")
	userID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	_, err := uuid.Parse(userID)
	if err != nil {
		er := errors.Wrap(err, "invalid user id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	_, err = time.Parse("2006-01-02", startDate)
	if err != nil {
		er := errors.Wrap(err, "invalid start date").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		er := errors.Wrap(err, "invalid end date").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := h.ExtraClient.TrackActivity(ctx, &pb.Period{
		Id:        userID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		er := errors.Wrap(err, "error tracking activity").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("TrackActivity method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// SetWorkingHours godoc
// @Summary Sets working hours
// @Description Sets working hours for kitchen
// @Tags kitchen
// @Param kitchen_id path string true "Kitchen ID"
// @Success 200 {object} extra.WorkingHoursResp
// @Failure 400 {object} string "Invalid kitchen ID or data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /local-eats/kitchens/{id}/working-hours [post]
func (h *Handler) SetWorkingHours(c *gin.Context) {
	h.Logger.Info("SetWorkingHours method is starting")
	kitchenID := c.Param("id")

	_, err := uuid.Parse(kitchenID)
	if err != nil {
		er := errors.Wrap(err, "invalid kitchen id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	var data map[string]*pb.DaySchedule
	if err := c.ShouldBindJSON(&data); err != nil {
		er := errors.Wrap(err, "invalid data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.ExtraClient.SetWorkingHours(ctx, &pb.WorkingHours{
		KitchenId: kitchenID,
		Schedule:  data,
	})
	if err != nil {
		er := errors.Wrap(err, "error setting working hours").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("SetWorkingHours method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// GetNutrition godoc
// @Summary Gets dish's nutrition info
// @Description Informs about dish's nutritional value
// @Tags dish
// @Param dish_id path string true "Dish ID"
// @Success 200 {object} extra.NutritionalInfo
// @Failure 400 {object} string "Invalid dish ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /local-eats/dishes/{id}/nutrition-info [get]
func (h *Handler) GetNutrition(c *gin.Context) {
	h.Logger.Info("GetNutrition method is starting")
	dishID := c.Param("id")

	_, err := uuid.Parse(dishID)
	if err != nil {
		er := errors.Wrap(err, "invalid dish id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.ExtraClient.GetNutrition(ctx, &pb.ID{Id: dishID})
	if err != nil {
		er := errors.Wrap(err, "error getting dish's nutritional info").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("GetNutrition method has finished successfully")
	c.JSON(http.StatusOK, res)
}
