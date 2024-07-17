package handler

import (
	pb "api-gateway/genproto/user"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// GetUser godoc
// @Summary Gets a user
// @Description Retrieves user info from database
// @Tags user
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} user.Profile
// @Failure 400 {object} string "Invalid user ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	h.Logger.Info("GetUser method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid user id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	profile, err := h.UserClient.GetProfile(ctx, &pb.ID{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error getting user").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("GetUser method has finished successfully")
	c.JSON(http.StatusOK, profile)
}

// UpdateUser godoc
// @Summary Updates a user
// @Description Updates user info in database
// @Tags user
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param user body user.NewInfoNoID true "User info"
// @Success 200 {object} user.Details
// @Failure 400 {object} string "Invalid user ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	h.Logger.Info("UpdateUser method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid user id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	var newData pb.NewInfoNoID
	if err := c.ShouldBindJSON(&newData); err != nil {
		er := errors.Wrap(err, "invalid user data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	upd, err := h.UserClient.UpdateProfile(ctx, &pb.NewInfo{
		Id:          id,
		FullName:    newData.FullName,
		Address:     newData.Address,
		PhoneNumber: newData.PhoneNumber,
	})
	if err != nil {
		er := errors.Wrap(err, "error updating user").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("UpdateUser method has finished successfully")
	c.JSON(http.StatusOK, upd)
}

// DeleteUser godoc
// @Summary Deletes a user
// @Description Deletes user from database
// @Tags user
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} user.Void
// @Failure 400 {object} string "Invalid user ID"
// @Failure 500 {object} string "Server error while processing request"
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	h.Logger.Info("DeleteUser method is starting")

	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		er := errors.Wrap(err, "invalid user id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	_, err = h.UserClient.DeleteProfile(ctx, &pb.ID{Id: id})
	if err != nil {
		er := errors.Wrap(err, "error deleting user").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er})
		h.Logger.Error(er)
		return
	}

	h.Logger.Info("DeleteUser method has finished successfully")
	c.JSON(http.StatusOK, "User deleted successfully")
}
