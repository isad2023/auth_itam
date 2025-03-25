package handlers

import (
	"fmt"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequestInput struct {
	Description string `json:"description" binding:"required"`
	Certificate string `json:"certificate"`
	Status      string `json:"status"`
	Type        string `json:"type" binding:"required"`
}

type UpdateRequestStatusRequest struct {
	RequestID uuid.UUID `json:"request_id" binding:"required"`
	Status    string    `json:"status" binding:"required"`
}

// @Summary Создать запрос пользователя
// @Description Создает новый запрос от имени пользователя
// @Tags Requests
// @Accept json
// @Produce json
// @Param request body handlers.CreateRequestInput true "Request data"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Success message with request ID"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/create_user_request [post]
func CreateUserRequest(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateRequestInput
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request := models.Request{
			ID:          uuid.New(),
			Description: input.Description,
			Certificate: input.Certificate,
			Status:      input.Status,
			Type:        input.Type,
			CreatedAt:   time.Now(),
		}

		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userObj, ok := user.(models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			return
		}
		request.UserID = userObj.ID

		if request.Status == "" {
			request.Status = "pending"
		}

		ctx := c.Request.Context()
		id, err := storage.SaveRequest(ctx, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Request created successfully", "request_id": id})
	}
}

// @Summary Получить запросы пользователя
// @Description Возвращает список запросов пользователя с пагинацией
// @Tags Requests
// @Produce json
// @Param user_id query string true "User ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Request data"
// @Failure 400 {object} map[string]string "Invalid user ID or pagination parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/get_request [get]
func GetRequest(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Query("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

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

		ctx := c.Request.Context()
		requests, err := storage.GetRequests(ctx, userID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": requests})
	}
}

// @Summary Получить все запросы пользователя
// @Description Возвращает список всех запросов текущего пользователя с пагинацией
// @Tags Requests
// @Produce json
// @Param user_id query string true "User ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "All requests"
// @Failure 400 {object} map[string]string "Invalid user ID or pagination parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/get_all_requests [get]
func GetAllRequests(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Query("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

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

		ctx := c.Request.Context()
		requests, err := storage.GetRequests(ctx, userID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": requests})
	}
}

// @Summary Обновить статус запроса
// @Description Обновляет статус указанного запроса
// @Tags Requests
// @Accept json
// @Produce json
// @Param request body handlers.UpdateRequestStatusRequest true "Request status update data"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/update_request_status [patch]
func UpdateRequestStatus(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			RequestID uuid.UUID `json:"request_id"`
			Status    string    `json:"status"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		err := storage.UpdateRequestStatus(ctx, request.RequestID, request.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "update_request_status"})
	}
}

// @Summary Удалить запрос
// @Description Удаляет запрос по его ID
// @Tags Requests
// @Produce json
// @Param request_id query string true "Request ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request ID"
// @Failure 404 {object} map[string]string "Request not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/delete_request [delete]
func DeleteRequest(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestIDParam := c.Query("request_id")
		if requestIDParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
			return
		}

		requestID, err := uuid.Parse(requestIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
			return
		}

		ctx := c.Request.Context()
		err = storage.DeleteRequest(ctx, requestID)
		if err != nil {
			if err.Error() == fmt.Sprintf("no request found with ID: %s", requestID) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting request", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Request deleted successfully"})
	}
}
