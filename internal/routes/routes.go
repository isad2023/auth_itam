package routes

import (
	"itam_auth/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes() *gin.Engine {

	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	config.AllowMethods = []string{"*"}
	router.Use(cors.New(config))

	// @Summary Пинг-сервис
	// @Description Проверяет доступность сервера
	// @Tags Health
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "pong"}
	// @Router /api/ping [get]
	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//* USER

	// @Summary Логин пользователя
	// @Description Авторизация пользователя с использованием логина и пароля
	// @Tags User
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "login"}
	// @Router /api/login [post]
	router.POST("/api/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "login",
		})
	})

	// @Summary Получить информацию о пользователе
	// @Description Возвращает данные текущего пользователя
	// @Tags User
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_user"}
	// @Router /api/get_user [get]
	router.GET("/api/get_user", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_user",
		})
	})

	// @Summary Обновить информацию пользователя
	// @Description Обновляет профиль пользователя
	// @Tags User
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "update_user_info"}
	// @Router /api/update_user_info [patch]
	router.PATCH("/api/update_user_info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "update_user_info",
		})
	})

	// @Summary Получить роли пользователя
	// @Description Возвращает список ролей текущего пользователя
	// @Tags User
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_user_roles"}
	// @Router /api/get_user_roles [get]
	router.GET("/api/get_user_roles", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_user_roles",
		})
	})

	// @Summary Получить свойства пользователя
	// @Description Возвращает список свойств текущего пользователя
	// @Tags User
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_user_properties"}
	// @Router /api/get_user_properties [get]
	router.GET("/api/get_user_properties", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_user_properties",
		})
	})

	// @Summary Получить достижения пользователя
	// @Description Возвращает список достижений текущего пользователя
	// @Tags User
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_user_achievements"}
	// @Router /api/get_user_achievements [get]
	router.GET("/api/get_user_achievements", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_user_achievements",
		})
	})

	//* Requests

	// @Summary Создать запрос пользователя
	// @Description Создает новый запрос от имени пользователя
	// @Tags Requests
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "create_user_request"}
	// @Router /api/create_user_request [post]
	router.POST("/api/create_user_request", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "create_user_request",
		})
	})

	// @Summary Получить запрос пользователя
	// @Description Возвращает данные о конкретном запросе пользователя
	// @Tags Requests
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_request"}
	// @Router /api/get_request [get]
	router.GET("/api/get_request", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_request",
		})
	})

	// @Summary Получить все запросы пользователя
	// @Description Возвращает список всех запросов текущего пользователя
	// @Tags Requests
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_all_requests"}
	// @Router /api/get_all_requests [get]
	router.GET("/api/get_all_requests", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_all_requests",
		})
	})

	// @Summary Обновить статус запроса
	// @Description Обновляет статус указанного запроса
	// @Tags Requests
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "update_request_status"}
	// @Router /api/update_request_status [patch]
	router.PATCH("/api/update_request_status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "update_request_status",
		})
	})

	//* Achievements

	// @Summary Создать достижение
	// @Description Создает новое достижение
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "create_achievement"}
	// @Router /api/create_achievement [post]
	router.POST("/api/create_achievement", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "create_achievement",
		})
	})

	// @Summary Обновить достижение
	// @Description Обновляет существующее достижение
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "update_achievement"}
	// @Router /api/update_achievement [patch]
	router.PATCH("/api/update_achievement", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "update_achievement",
		})
	})

	// @Summary Получить достижение
	// @Description Возвращает информацию о конкретном достижении
	// @Tags Achievements
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_achievement"}
	// @Router /api/get_achievement [get]
	router.GET("/api/get_achievement", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_achievement",
		})
	})

	// @Summary Получить все достижения
	// @Description Возвращает список всех достижений
	// @Tags Achievements
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_all_achievements"}
	// @Router /api/get_all_achievements [get]
	router.GET("/api/get_all_achievements", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_all_achievements",
		})
	})

	//* Notifications (admin)

	// CreateNotification godoc
	// @Summary Create a new notification
	// @Description Create a new notification
	// @Tags Notifications
	// @Accept json
	// @Produce json
	// @Param notification body models.Notification true "Notification"
	// @Success 201 {object} gin.H{"message": "Notification created successfully"}
	// @Failure 400 {object} gin.H{"error": "Invalid request"}
	// @Failure 500 {object} gin.H{"error": "Error while saving notification"}
	// @Router /api/create_notification [post]
	router.POST("/api/create_notification", handlers.CreateNotification)

	// UpdateNotification godoc
	// @Summary Update an existing notification
	// @Description Update an existing notification
	// @Tags Notifications
	// @Accept json
	// @Produce json
	// @Param notification body models.Notification true "Notification"
	// @Success 201 {object} gin.H{"message": "Notification updated successfully"}
	// @Failure 400 {object} gin.H{"error": "Invalid request"}
	// @Failure 500 {object} gin.H{"error": "Error while updating notification"}
	// @Router /api/update_notification [patch]
	router.PATCH("/api/update_notification", handlers.UpdateNotification)

	// GetAllNotifications godoc
	// @Summary Get all notifications for a user
	// @Description Get all notifications for a user
	// @Tags Notifications
	// @Produce json
	// @Param user_id path string true "User ID"
	// @Success 201 {object} models.Notification
	// @Failure 400 {object} gin.H{"error": "Invalid user ID"}
	// @Failure 500 {object} gin.H{"error": "Error while fetching notifications"}
	// @Router /api/get_all_notifications/{user_id} [get]
	router.GET("/api/get_all_notifications", handlers.GetAllNotifications)

	// GetNotification godoc
	// @Summary Get a notification by its ID
	// @Description Get a notification by its ID
	// @Tags Notifications
	// @Produce json
	// @Param notification_id path string true "Notification ID"
	// @Success 201 {object} models.Notification
	// @Failure 400 {object} gin.H{"error": "Invalid notification ID"}
	// @Failure 500 {object} gin.H{"error": "Error while fetching notification"}
	// @Router /api/get_notification/{notification_id} [get]
	router.GET("/api/get_notification", handlers.GetNotification)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
