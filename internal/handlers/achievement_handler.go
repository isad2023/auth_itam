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

func CreateAchievement(c *gin.Context, storage *database.Storage) {
	var achievement models.Achievement
	if err := c.BindJSON(&achievement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Println("Received achievement:", achievement)

	achievement.ID = uuid.New()
	achievement.CreatedAt = time.Now()

	ctx := context.Background()

	log.Printf("Saving achievement: %+v\n", achievement)
	_, err := storage.SaveAchievement(ctx, achievement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving achievement: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Achievement created successfully"})
}

func UpdateAchievement(c *gin.Context, storage *database.Storage) {
	var achievement models.Achievement
	if err := c.BindJSON(&achievement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()
	err := storage.UpdateAchievement(ctx, achievement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating achievement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Achievement updated successfully"})
}

func GetAllAchievements(c *gin.Context, storage *database.Storage) {
	ctx := context.Background()
	achievements, err := storage.GetAllAchievements(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

func GetAchievementByID(c *gin.Context, storage *database.Storage) {
	achievementID := c.Query("achievement_id")
	uuidAchievementID, err := uuid.Parse(achievementID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement ID"})
		return
	}

	ctx := context.Background()
	achievement, err := storage.GetAchievementByID(ctx, uuidAchievementID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching achievement"})
		return
	}

	c.JSON(http.StatusOK, achievement)
}

func GetAchievementsByUserID(c *gin.Context, storage *database.Storage) {
	userID := c.Query("user_id")
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()
	achievements, err := storage.GetAchievementsByUserID(ctx, uuidUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

func DeleteAchievement(c *gin.Context, storage *database.Storage) {
	achievementID := c.Query("achievement_id")
	uuidAchievementID, err := uuid.Parse(achievementID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement ID"})
		return
	}

	ctx := context.Background()
	err = storage.DeleteAchievement(ctx, uuidAchievementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting achievement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Achievement deleted successfully"})
}
