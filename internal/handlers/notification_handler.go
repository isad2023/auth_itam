package handlers

import (
	"context"
	"database/sql"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateNotification(c *gin.Context, storage *database.Storage) {
	var notification models.Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Println("Received notification:", notification)

	notification.ID = uuid.New()
	notification.CreatedAt = time.Now()
	notification.IsRead = false

	ctx := context.Background()

	log.Printf("Saving notification: %+v\n", notification)
	_, err := storage.SaveNotification(ctx, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving notification:" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notification created successfully"})
}

func UpdateNotification(c *gin.Context, storage *database.Storage) {
	var notification models.Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()
	err := storage.UpdateNotification(ctx, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification updated successfully"})
}

func GetAllNotifications(c *gin.Context, storage *database.Storage) {
	userID := c.Query("user_id")
	var notifications []models.Notification
	var err error

	ctx := context.Background()
	if userID == "" {
		notifications, err = storage.GetAllNotifications(ctx)
	} else {
		uuidUserID, errUUID := uuid.Parse(userID)
		if errUUID != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		notifications, err = storage.GetNotifications(ctx, uuidUserID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func GetNotification(c *gin.Context, storage *database.Storage) {
	notificationID := c.Param("notification_id")
	uuidNotificationID, err := uuid.Parse(notificationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	ctx := context.Background()
	notification, err := storage.GetNotificationByID(ctx, uuidNotificationID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching notification"})
		return
	}

	c.JSON(http.StatusOK, notification)
}
