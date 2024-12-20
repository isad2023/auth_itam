package handlers

import (
	"context"
	"database/sql"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	pgUUID "github.com/jackc/pgtype/ext/gofrs-uuid"
)

func CreateNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	notificationID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating notification ID"})
		return
	}
	notification.ID = pgUUID.UUID{UUID: notificationID}
	notification.CreatedAt = time.Now()
	notification.IsRead = false

	ctx := context.Background()
	_, err = database.SaveNotification(ctx, database.DB, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notification created successfully"})
}

func UpdateNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()
	err := database.UpdateNotification(ctx, database.DB, notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification updated successfully"})
}

func GetAllNotifications(c *gin.Context) {
	userID := c.Param("user_id")
	uuidUserID, err := uuid.FromString(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()
	notifications, err := database.GetNotifications(ctx, database.DB, uuidUserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func GetNotification(c *gin.Context) {
	notificationID := c.Param("notification_id")
	uuidNotificationID, err := uuid.FromString(notificationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	ctx := context.Background()
	notification, err := database.GetNotificationByID(ctx, database.DB, uuidNotificationID)

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
