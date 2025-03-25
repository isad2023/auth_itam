package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func parseIntQuery(c *gin.Context, paramName string, defaultValue int) (int, error) {
	param := c.Query(paramName)
	if param == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(param)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("invalid %s, must be a non-negative integer", paramName)
	}

	return value, nil
}

// @Summary Создать достижение
// @Description Создает новое достижение
// @Tags Achievements
// @Accept json
// @Produce json
// @Param achievement body models.Achievement true "Achievement data"
// @Security OAuth2Password
// @Success 201 {object} map[string]interface{} "Success message with ID"
// @Failure 400 {object} map[string]string "Invalid title or points"
// @Failure 500 {object} map[string]string "Failed to save achievement"
// @Router /api/create_achievement [post]
func CreateAchievement(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var achievement models.Achievement
		if err := c.BindJSON(&achievement); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if achievement.Title == "" || achievement.Points < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid title or points"})
			return
		}

		achievement.ID = uuid.New()
		achievement.CreatedAt = time.Now()

		ctx := c.Request.Context()
		if _, err := storage.SaveAchievement(ctx, achievement); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save achievement"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Achievement created successfully", "id": achievement.ID})
	}
}

// @Summary Обновить достижение
// @Description Обновляет существующее достижение
// @Tags Achievements
// @Accept json
// @Produce json
// @Param achievement body models.Achievement true "Achievement data"
// @Security OAuth2Password
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid achievement ID or data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/update_achievement [patch]
func UpdateAchievement(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var achievement models.Achievement
		if err := c.BindJSON(&achievement); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if achievement.ID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement ID"})
			return
		}

		ctx := c.Request.Context()
		if err := storage.UpdateAchievement(ctx, achievement); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievement"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Achievement updated successfully"})
	}
}

// @Summary Получить все достижения
// @Description Возвращает список всех достижений с пагинацией
// @Tags Achievements
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Security OAuth2Password
// @Success 200 {array} models.Achievement "List of all achievements"
// @Failure 400 {object} map[string]string "Invalid pagination parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/get_all_achievements [get]
func GetAllAchievements(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := parseIntQuery(c, "limit", 10)
		if err != nil {
			return
		}

		offset, err := parseIntQuery(c, "offset", 0)
		if err != nil {
			return
		}

		ctx := c.Request.Context()
		achievements, err := storage.GetAllAchievements(ctx, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
			return
		}

		c.JSON(http.StatusOK, achievements)
	}
}

// @Summary Получить достижение
// @Description Возвращает информацию о конкретном достижении
// @Tags Achievements
// @Produce json
// @Param achievement_id query string true "Achievement ID"
// @Security OAuth2Password
// @Success 200 {object} models.Achievement "Achievement data"
// @Failure 400 {object} map[string]string "Invalid achievement ID"
// @Failure 404 {object} map[string]string "Achievement not found"
// @Failure 500 {object} map[string]string "Error while fetching achievement"
// @Router /api/get_achievement [get]
func GetAchievementByID(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("achievement_id")
		achievementID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement ID"})
			return
		}

		ctx := c.Request.Context()
		achievement, err := storage.GetAchievementByID(ctx, achievementID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching achievement"})
			return
		}

		c.JSON(http.StatusOK, achievement)
	}
}

// @Summary Получить достижения пользователя
// @Description Возвращает список достижений пользователя с пагинацией
// @Tags Achievements
// @Produce json
// @Param user_id query string true "User ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Security OAuth2Password
// @Success 200 {array} models.Achievement "List of achievements"
// @Failure 400 {object} map[string]string "Invalid user ID or pagination parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/get_user_achievements [get]
func GetAchievementsByUserID(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("user_id")
		userID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		limit, err := parseIntQuery(c, "limit", 10)
		if err != nil {
			return
		}

		offset, err := parseIntQuery(c, "offset", 0)
		if err != nil {
			return
		}

		ctx := c.Request.Context()
		achievements, err := storage.GetAchievementsByUserID(ctx, userID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching achievements"})
			return
		}

		c.JSON(http.StatusOK, achievements)
	}
}

// @Summary Удалить достижение
// @Description Удаляет достижение по его ID
// @Tags Achievements
// @Produce json
// @Param achievement_id query string true "Achievement ID"
// @Security OAuth2Password
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid achievement ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/delete_achievement [delete]
func DeleteAchievement(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("achievement_id")
		achievementID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement ID"})
			return
		}

		ctx := c.Request.Context()
		err = storage.DeleteAchievement(ctx, achievementID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting achievement"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Achievement deleted successfully"})
	}
}
