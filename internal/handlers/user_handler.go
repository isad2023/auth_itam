package handlers

import (
	"context"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"itam_auth/internal/services/auth"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// type TelegramAuth struct {
// 	ID        int64  `json:"id"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Username  string `json:"username"`
// 	PhotoURL  string `json:"photo_url"`
// 	AuthDate  int64  `json:"auth_date"`
// 	Hash      string `json:"hash"`
// }

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"username" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

// @Summary Регистрация нового пользователя
// @Description Регистрация нового пользователя в системе
// @Tags User
// @Accept json
// @Produce json
// @Param register body handlers.RegisterRequest true "User registration details"
// @Success 201 {object} map[string]string "Success message with user data"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/register [post]
func Register(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		user, err := auth.RegisterUser(ctx, storage, req.Name, req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving user", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
	}
}

// @Summary Логин пользователя
// @Description Авторизация пользователя с использованием логина и пароля
// @Tags User
// @Accept json,x-www-form-urlencoded
// @Produce json
// @Param login body handlers.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "JWT token"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/login [post]
func Login(storage *database.Storage, hmacSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		
		// Поддержка как JSON, так и form-data
		if c.ContentType() == "application/json" {
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		} else {
			if err := c.ShouldBind(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		ctx := context.Background()
		tokenString, err := auth.AuthenticateUser(ctx, storage, req.Email, req.Password, hmacSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password", "details": err.Error()})
			return
		}

		// Формат ответа для совместимости с OAuth2
		c.JSON(http.StatusOK, gin.H{
			"access_token": tokenString,
			"token_type": "Bearer",
			"expires_in": 2592000, // 30 дней в секундах
		})
	}
}

// TELEGRAM WEB APP
// func Login(c *gin.Context, storage *database.Storage) {
// 	var auth TelegramAuth
// 	if err := c.ShouldBindJSON(&auth); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	data := url.Values{
// 		"id":         {fmt.Sprintf("%d", auth.ID)},
// 		"first_name": {auth.FirstName},
// 		"last_name":  {auth.LastName},
// 		"username":   {auth.Username},
// 		"photo_url":  {auth.PhotoURL},
// 		"auth_date":  {fmt.Sprintf("%d", auth.AuthDate)},
// 	}

// 	if !utils.ValidateTelegramAuth(data, "") {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Telegram authorization"})
// 		return
// 	}

// 	// Генерация UUID
// 	userID := uuid.New()

// 	ctx := context.Background()
// 	user, err := storage.GetUserByID(ctx, userID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	tokenString, err := utils.GenerateJWT(user.Email)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": tokenString})
// }

// TELEGRAM WEB APP
// func Register(c *gin.Context, storage *database.Storage) {
// 	var auth TelegramAuth

// 	if err := c.ShouldBindJSON(&auth); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	data := url.Values{
// 		"id":         {fmt.Sprintf("%d", auth.ID)},
// 		"first_name": {auth.FirstName},
// 		"last_name":  {auth.LastName},
// 		"username":   {auth.Username},
// 		"photo_url":  {auth.PhotoURL},
// 		"auth_date":  {fmt.Sprintf("%d", auth.AuthDate)},
// 	}

// 	if !utils.ValidateTelegramAuth(data, "your_bot_token") {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Telegram authorization"})
// 		return
// 	}

// 	userID := uuid.New()

// 	user := models.User{
// 		ID:           userID,
// 		Name:         fmt.Sprintf("%s %s", auth.FirstName, auth.LastName),
// 		Email:        fmt.Sprintf("%d@telegram.com", auth.ID),
// 		Telegram:     auth.Username,
// 		PasswordHash: "",
// 		PhotoURL:     auth.PhotoURL,
// 		CreatedAt:    time.Now(),
// 		UpdatedAt:    time.Now(),
// 	}

// 	ctx := context.Background()

// 	_, err := storage.SaveUser(ctx, user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
// }

// @Summary Получить информацию о пользователе
// @Description Возвращает данные текущего пользователя
// @Tags User
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.User "User data"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/get_user/{user_id} [get]
func GetUser(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		uuidUserID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		ctx := context.Background()
		user, err := storage.GetUserByID(ctx, uuidUserID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// @Summary Получить роли пользователя
// @Description Возвращает список ролей текущего пользователя
// @Tags User
// @Produce json
// @Security OAuth2Password
// @Success 200 {array} models.UserRole "User roles"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/get_user_roles [get]
func GetUserRoles(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			log.Printf("RequestID: %s, Error: User not authenticated, Details: User not found in context", c.GetString("request_id"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated", "details": "User not found in context"})
			return
		}
		userObj, ok := user.(models.User)
		if !ok {
			log.Printf("RequestID: %s, Error: Invalid user data, Details: User object is invalid", c.GetString("request_id"))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data", "details": "User object is invalid"})
			return
		}
		userID := userObj.ID

		ctx := c.Request.Context()
		roles, err := storage.GetUserRoles(ctx, userID)
		if err != nil {
			log.Printf("RequestID: %s, Error: Error while fetching user roles, Details: %v", c.GetString("request_id"), err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching user roles", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, roles)
	}
}

// @Summary Получить свойства пользователя
// @Description Возвращает список свойств текущего пользователя
// @Tags User
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} map[string]string "User properties"
// @Router /api/get_user_properties [get]
func GetUserPermissions(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user from context that was set by AuthMiddleware
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
		
		uuidUserID := userObj.ID

		ctx := context.Background()
		userRoles, err := storage.GetUserRoles(ctx, uuidUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching user roles", "details": err.Error()})
			return
		}

		var permissions []models.RolePermission
		for _, role := range userRoles {
			rolePermissions, err := storage.GetRolePermissions(ctx, role.RoleID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching role permissions", "details": err.Error()})
				return
			}
			permissions = append(permissions, rolePermissions...)
		}

		c.JSON(http.StatusOK, permissions)
	}
}

// @Summary Обновить информацию пользователя
// @Description Обновляет профиль пользователя
// @Tags User
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param user body models.User true "User update data"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/update_user_info [patch]
func UpdateUserInfo(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated user from context
		userAuth, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		
		authUser, ok := userAuth.(models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			return
		}
		
		// Get update data from request
		var updateData models.User
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}
		
		// Only allow updating certain fields
		user := models.User{
			ID:            authUser.ID,
			Name:          updateData.Name,
			Specification: updateData.Specification,
			About:         updateData.About,
			PhotoURL:      updateData.PhotoURL,
			ResumeURL:     updateData.ResumeURL,
			UpdatedAt:     time.Now(),
		}
		
		// Update user in database
		ctx := context.Background()
		err := storage.UpdateUser(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "User information updated successfully"})
	}
}

// @Summary Получить информацию о текущем пользователе
// @Description Возвращает данные авторизованного пользователя
// @Tags User
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} models.User "User data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/me [get]
func GetCurrentUser(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user from context that was set by AuthMiddleware
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
		
		// Get full user information from database
		ctx := context.Background()
		fullUser, err := storage.GetUserByID(ctx, userObj.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user data", "details": err.Error()})
			return
		}
		
		// Remove sensitive information
		fullUser.PasswordHash = ""
		
		c.JSON(http.StatusOK, fullUser)
	}
}
