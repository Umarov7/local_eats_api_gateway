package handler

import (
	pb "api-gateway/genproto/dish"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateDish godoc
// @Summary Creates a dish
// @Description Inserts a new dish into database
// @Tags dish
// @Security ApiKeyAuth
// @Param dish body dish.NewDish true "Dish info"
// @Success 200 {object} dish.NewDishResp
// @Failure 400 {object} string "Invalid dish data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /dishes [post]
func (h *Handler) CreateDish(c *gin.Context) {
	h.Logger.Info("CreateDish method is starting")

	var data pb.NewDish
	if err := c.ShouldBindJSON(&data); err != nil {
		er := errors.Wrap(err, "invalid dish data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.DishClient.Add(ctx, &data)
	if err != nil {
		er := errors.Wrap(err, "error creating dish").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("CreateDish method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// GetDish godoc
// @Summary Gets a dish
// @Description Retrieves dish info from database
// @Tags dish
// @Security ApiKeyAuth
// @Param id path string true "Dish ID"
// @Success 200 {object} dish.DishInfo
// @Failure 400 {object} string "Invalid dish ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /dishes/{id} [get]
func (h *Handler) GetDish(c *gin.Context) {
	h.Logger.Info("GetDish method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid dish ID").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.DishClient.Read(ctx, &pb.ID{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting dish").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("GetDish method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// UpdateDish godoc
// @Summary Updates a dish
// @Description Updates dish info in database
// @Tags dish
// @Security ApiKeyAuth
// @Param id path string true "Dish ID"
// @Param dish body dish.NewDataNoID true "Dish info"
// @Success 200 {object} dish.UpdatedData
// @Failure 400 {object} string "Invalid dish ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /dishes/{id} [put]
func (h *Handler) UpdateDish(c *gin.Context) {
	h.Logger.Info("UpdateDish method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid dish ID").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	var data pb.NewData
	if err := c.ShouldBindJSON(&data); err != nil {
		er := errors.Wrap(err, "invalid dish data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.DishClient.Update(ctx, &pb.NewData{
		Id:        id,
		Name:      data.Name,
		Price:     data.Price,
		Available: data.Available,
	})
	if err != nil {
		er := errors.Wrap(err, "error updating dish").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("UpdateDish method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// DeleteDish godoc
// @Summary Deletes a dish
// @Description Deletes dish from database
// @Tags dish
// @Security ApiKeyAuth
// @Param id path string true "Dish ID"
// @Success 200 {object} string
// @Failure 400 {object} string "Invalid dish ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /dishes/{id} [delete]
func (h *Handler) DeleteDish(c *gin.Context) {
	h.Logger.Info("DeleteDish method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid dish ID").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	_, err = h.DishClient.Delete(ctx, &pb.ID{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error deleting dish").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("DeleteDish method has finished successfully")
	c.JSON(http.StatusOK, "Dish deleted successfully")
}

// FetchDishes godoc
// @Summary Gets dishes
// @Description Retrieves dishes info from database
// @Tags dish
// @Security ApiKeyAuth
// @Param id path string true "Kitchen ID"
// @Param page query int true "Page number"
// @Param limit query int true "Number of items per page"
// @Success 200 {object} dish.Dishes
// @Failure 500 {object} string "Server error while processing request"
// @Router /kitchens/{id}/dishes [get]
func (h *Handler) FetchDishes(c *gin.Context) {
	h.Logger.Info("FetchDishes method is starting")

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

	res, err := h.DishClient.Fetch(ctx, &pb.Pagination{
		Limit:  int32(l),
		Offset: int32((p - 1) * l),
	})
	if err != nil {
		er := errors.Wrap(err, "error getting dishes").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("FetchDishes method has finished successfully")
	c.JSON(http.StatusOK, res)
}
