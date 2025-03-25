package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Создать уведомление
// @Description Создает новое уведомление
// @Tags Notifications
// @Accept json
// @Produce json
// @Param notification body models.Notification true "Notification data"
// @Security BearerAuth
// @Success 201 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/create_notification [post]
func CreateNotification(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notification models.Notification
		if err := c.ShouldBindJSON(&notification); err != nil {
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
}

// UpdateNotification обновляет существующее уведомление
// @Summary Обновить уведомление
// @Description Обновляет существующее уведомление
// @Tags Notifications
// @Accept json
// @Produce json
// @Param notification body models.Notification true "Notification data"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/update_notification [patch]
func UpdateNotification(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notification models.Notification
		if err := c.ShouldBindJSON(&notification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if notification.ID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
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
}

// @Summary Получить все уведомления
// @Description Возвращает список всех уведомлений или уведомлений пользователя с пагинацией
// @Tags Notifications
// @Produce json
// @Param user_id query string false "User ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Security BearerAuth
// @Success 200 {array} models.Notification "List of notifications"
// @Failure 400 {object} map[string]string "Invalid user ID or pagination parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/get_all_notifications [get]
func GetAllNotifications(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := parseIntQuery(c, "limit", 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
			return
		}

		offset, err := parseIntQuery(c, "offset", 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
			return
		}

		userID := c.Query("user_id")
		ctx := context.Background()
		var notifications []models.Notification

		if userID == "" {
			notifications, err = storage.GetAllNotifications(ctx, limit, offset)
		} else {
			uuidUserID, errUUID := uuid.Parse(userID)
			if errUUID != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}
			notifications, err = storage.GetNotifications(ctx, uuidUserID, limit, offset)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching notifications"})
			return
		}

		c.JSON(http.StatusOK, notifications)
	}
}

// @Summary Получить уведомление по ID
// @Description Возвращает уведомление по его ID
// @Tags Notifications
// @Produce json
// @Param notification_id path string true "Notification ID"
// @Security BearerAuth
// @Success 200 {object} models.Notification "Notification data"
// @Failure 400 {object} map[string]string "Invalid notification ID"
// @Failure 404 {object} map[string]string "Notification not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/get_notification/{notification_id} [get]
func GetNotification(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

// @Summary Удалить уведомление
// @Description Удаляет уведомление по ID
// @Tags Notifications
// @Produce json
// @Param notification_id query string true "Notification ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid notification ID"
// @Failure 404 {object} map[string]string "Notification not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/delete_notification [delete]
func DeleteNotification(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		notificationID := c.Query("notification_id")
		if notificationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
			return
		}

		uuidNotificationID, err := uuid.Parse(notificationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
			return
		}

		ctx := context.Background()
		err = storage.DeleteNotification(ctx, uuidNotificationID)
		if err != nil {
			if err.Error() == fmt.Sprintf("no notification found with ID: %s", uuidNotificationID) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting notification", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
	}
}
