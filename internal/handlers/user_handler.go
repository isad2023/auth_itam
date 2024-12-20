package handlers

import (
	"context"
	"fmt"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"itam_auth/internal/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	gofrs_uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type TelegramAuth struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

func Login(c *gin.Context, storage *database.Storage) {
	var auth TelegramAuth
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	data := url.Values{
		"id":         {fmt.Sprintf("%d", auth.ID)},
		"first_name": {auth.FirstName},
		"last_name":  {auth.LastName},
		"username":   {auth.Username},
		"photo_url":  {auth.PhotoURL},
		"auth_date":  {fmt.Sprintf("%d", auth.AuthDate)},
	}

	if !utils.ValidateTelegramAuth(data, "") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Telegram authorization"})
		return
	}

	// Генерация UUID
	userID, _ := uuid.NewV4()

	ctx := context.Background()
	user, err := storage.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	tokenString, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Register(c *gin.Context, storage *database.Storage) {
	var auth TelegramAuth

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	data := url.Values{
		"id":         {fmt.Sprintf("%d", auth.ID)},
		"first_name": {auth.FirstName},
		"last_name":  {auth.LastName},
		"username":   {auth.Username},
		"photo_url":  {auth.PhotoURL},
		"auth_date":  {fmt.Sprintf("%d", auth.AuthDate)},
	}

	if !utils.ValidateTelegramAuth(data, "your_bot_token") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Telegram authorization"})
		return
	}

	userID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating user ID"})
		return
	}

	pgUUID := gofrs_uuid.UUID{UUID: userID}

	user := models.User{
		ID:           pgUUID,
		Name:         fmt.Sprintf("%s %s", auth.FirstName, auth.LastName),
		Email:        fmt.Sprintf("%d@telegram.com", auth.ID),
		Telegram:     auth.Username,
		PasswordHash: "",
		PhotoURL:     auth.PhotoURL,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	ctx := context.Background()

	_, err = storage.SaveUser(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func GetUser(c *gin.Context, storage *database.Storage) {
	userID := c.Param("user_id")
	uuidUserID, err := uuid.FromString(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()
	user, err := storage.GetUserByID(ctx, uuidUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
