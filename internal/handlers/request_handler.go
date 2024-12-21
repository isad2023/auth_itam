package handlers

import (
	"context"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func CreateUserRequest(c *gin.Context, storage *database.Storage) {
	var request models.Request

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.ID, _ = uuid.NewV4()
	request.UserID, _ = uuid.NewV4()
	request.CreatedAt = time.Now()

	ctx := context.Background()
	id, err := storage.SaveRequest(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "create_user_request", "request_id": id})
}

func GetRequest(c *gin.Context, storage *database.Storage) {
	userID, err := uuid.FromString(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()
	requests, err := storage.GetRequests(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": requests})
}

func GetAllRequests(c *gin.Context, storage *database.Storage) {
	userID, err := uuid.FromString(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()
	requests, err := storage.GetRequests(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": requests})
}

func UpdateRequestStatus(c *gin.Context, storage *database.Storage) {
	var request struct {
		RequestID uuid.UUID `json:"request_id"`
		Status    string    `json:"status"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err := storage.UpdateRequestStatus(ctx, request.RequestID, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update_request_status"})
}
