package handlers

import (
	"database/sql"
	"errors"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateAchievement(c *gin.Context, storage *database.Storage) {
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

func UpdateAchievement(c *gin.Context, storage *database.Storage) {
	var achievement models.Achievement
	if err := c.BindJSON(&achievement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := c.Request.Context()
	if err := storage.UpdateAchievement(ctx, achievement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Achievement updated successfully"})
}

func GetAllAchievements(c *gin.Context, storage *database.Storage) {
	ctx := c.Request.Context()
	achievements, err := storage.GetAllAchievements(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

func GetAchievementByID(c *gin.Context, storage *database.Storage) {
	idParam := c.Query("achievement_id")
	achievementID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement ID"})
		return
	}

	ctx := c.Request.Context()
	achievement, err := storage.GetAchievementByID(ctx, achievementID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // Проверяем обёрнутую ошибку
			c.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching achievement"})
		return
	}

	c.JSON(http.StatusOK, achievement)
}

func GetAchievementsByUserID(c *gin.Context, storage *database.Storage) {
	idParam := c.Query("user_id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := c.Request.Context()
	achievements, err := storage.GetAchievementsByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

func DeleteAchievement(c *gin.Context, storage *database.Storage) {
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
