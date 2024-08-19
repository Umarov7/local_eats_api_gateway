package handler

import (
	pb "api-gateway/genproto/kitchen"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateKitchen godoc
// @Summary Creates a kitchen
// @Description Inserts a new kitchen into database
// @Tags kitchen
// @Security ApiKeyAuth
// @Param kitchen body kitchen.CreateRequest true "Kitchen info"
// @Success 200 {object} kitchen.CreateResponse
// @Failure 400 {object} string "Invalid kitchen data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /kitchens [post]
func (h *Handler) CreateKitchen(c *gin.Context) {
	h.Logger.Info("CreateKitchen method is starting")
	var data pb.CreateRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		er := errors.Wrap(err, "invalid kitchen data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.KitchenClient.Create(ctx, &data)
	if err != nil {
		er := errors.Wrap(err, "error creating kitchen").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("CreateKitchen method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// GetKitchen godoc
// @Summary Gets a kitchen
// @Description Retrieves kitchen info from database
// @Tags kitchen
// @Security ApiKeyAuth
// @Param id path string true "Kitchen ID"
// @Success 200 {object} kitchen.Info
// @Failure 400 {object} string "Invalid kitchen ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /kitchens/{id} [get]
func (h *Handler) GetKitchen(c *gin.Context) {
	h.Logger.Info("GetKitchen method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid kitchen id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	kitchen, err := h.KitchenClient.Get(ctx, &pb.ID{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting kitchen").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("GetKitchen method has finished successfully")
	c.JSON(http.StatusOK, kitchen)
}

// UpdateKitchen godoc
// @Summary Updates a kitchen
// @Description Updates kitchen info in database
// @Tags kitchen
// @Security ApiKeyAuth
// @Param id path string true "Kitchen ID"
// @Param kitchen body kitchen.NewDataNoID true "Kitchen info"
// @Success 200 {object} kitchen.UpdatedData
// @Failure 400 {object} string "Invalid kitchen ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /kitchens/{id} [put]
func (h *Handler) UpdateKitchen(c *gin.Context) {
	h.Logger.Info("UpdateKitchen method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid kitchen id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	var data pb.NewDataNoID
	if err := c.ShouldBindJSON(&data); err != nil {
		er := errors.Wrap(err, "invalid kitchen data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	res, err := h.KitchenClient.Update(ctx, &pb.NewData{
		Id:          id,
		Name:        data.Name,
		Description: data.Description,
		PhoneNumber: data.PhoneNumber,
	})
	if err != nil {
		er := errors.Wrap(err, "error updating kitchen").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("UpdateKitchen method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// DeleteKitchen godoc
// @Summary Deletes a kitchen
// @Description Deletes kitchen from database
// @Tags kitchen
// @Security ApiKeyAuth
// @Param id path string true "Kitchen ID"
// @Success 200 {object} string
// @Failure 400 {object} string "Invalid kitchen ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /kitchens/{id} [delete]
func (h *Handler) DeleteKitchen(c *gin.Context) {
	h.Logger.Info("DeleteKitchen method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid kitchen id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	_, err = h.KitchenClient.Delete(ctx, &pb.ID{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error deleting kitchen").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("DeleteKitchen method has finished successfully")
	c.JSON(http.StatusOK, "Kitchen deleted successfully")
}

// FetchKitchens godoc
// @Summary Fetches all kitchens
// @Description Fetches all kitchens from database
// @Tags kitchen
// @Security ApiKeyAuth
// @Param page query int true "Page number"
// @Param limit query int true "Number of items per page"
// @Success 200 {object} kitchen.Kitchens
// @Failure 500 {object} string "Server error while processing request"
// @Router /kitchens [get]
func (h *Handler) FetchKitchens(c *gin.Context) {
	h.Logger.Info("FetchKitchens method is starting")

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

	res, err := h.KitchenClient.Fetch(ctx, &pb.Pagination{
		Limit:  int32(l),
		Offset: int32((p - 1) * l),
	})
	if err != nil {
		er := errors.Wrap(err, "error fetching kitchens").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("FetchKitchens method has finished successfully")
	c.JSON(http.StatusOK, res)
}

// SearchKitchens godoc
// @Summary Searches kitchens
// @Description Searches kitchens from database
// @Tags kitchen
// @Security ApiKeyAuth
// @Param query query string false "Search query"
// @Param cuisine_type query string false "Cuisine type"
// @Param rating query float32 false "Rating"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} kitchen.Kitchens
// @Failure 500 {object} string "Server error while processing request"
// @Router /kitchens/search [get]
func (h *Handler) SearchKitchens(c *gin.Context) {
	h.Logger.Info("SearchKitchens method is starting")

	query := c.Query("query")
	cuisineType := c.Query("cuisine_type")
	rating := c.Query("rating")
	page := c.Query("page")
	limit := c.Query("limit")
	var ratingFloat float64

	if query == "" && cuisineType == "" && rating == "" {
		er := errors.New("invalid search parameters").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	if rating != "" {
		r, err := strconv.ParseFloat(rating, 32)
		if err != nil {
			er := errors.Wrap(err, "invalid search parameters").Error()
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": er})
			h.Logger.Error(er)
			return
		}
		ratingFloat = r
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

	res, err := h.KitchenClient.Search(ctx, &pb.SearchDetails{
		Query:       query,
		CuisineType: cuisineType,
		Rating:      float32(ratingFloat),
		Pagination: &pb.Pagination{
			Limit:  int32(l),
			Offset: int32((p - 1) * l),
		},
	})
	if err != nil {
		er := errors.Wrap(err, "error searching kitchens").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("SearchKitchens method has finished successfully")
	c.JSON(http.StatusOK, res)
}
