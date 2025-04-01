package routes

import (
	"itam_auth/internal/database"
	"itam_auth/internal/handlers"
	"itam_auth/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(storage *database.Storage, hmacSecret string) *gin.Engine {

	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.Config{
		AllowOrigins:     []string{
			"http://localhost:3000",
			"http://45.10.41.58:3000",
			"http://localhost:8080",
			"http://45.10.41.58:8080",
			"http://localhost:8080/swagger",
			"http://45.10.41.58:8080/swagger",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}
	router.Use(cors.New(config))

	api := router.Group("/api")
	{
		// Public routes that don't require authorization
		api.GET("/ping", pingHandler)
		api.POST("/login", handlers.Login(storage, hmacSecret))
		api.POST("/register", handlers.Register(storage))
		api.GET("/get_user/:user_id", handlers.GetUser(storage))

		// Protected routes that require authorization
		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware(hmacSecret))
		{
			auth.GET("/me", handlers.GetCurrentUser(storage))
			auth.PATCH("/update_user_info", handlers.UpdateUserInfo(storage))
			auth.GET("/get_user_roles", handlers.GetUserRoles(storage))
			auth.GET("/get_user_properties", handlers.GetUserPermissions(storage))

			//* REQUEST ROUTES
			auth.POST("/create_user_request", handlers.CreateUserRequest(storage))
			auth.GET("/get_request", handlers.GetRequest(storage))
			auth.GET("/get_all_requests", handlers.GetAllRequests(storage))
			auth.PATCH("/update_request_status", handlers.UpdateRequestStatus(storage))
			auth.DELETE("/delete_request", handlers.DeleteRequest(storage))

			//* ACHIEVEMENT ROUTES
			auth.GET("/get_user_achievements", handlers.GetAchievementsByUserID(storage))
			auth.POST("/create_achievement", handlers.CreateAchievement(storage))
			auth.PATCH("/update_achievement", handlers.UpdateAchievement(storage))
			auth.GET("/get_achievement", handlers.GetAchievementByID(storage))
			auth.GET("/get_all_achievements", handlers.GetAllAchievements(storage))
			auth.DELETE("/delete_achievement", handlers.DeleteAchievement(storage))

			//* NOTIFICATION ROUTES
			auth.POST("/create_notification", handlers.CreateNotification(storage))
			auth.PATCH("/update_notification", handlers.UpdateNotification(storage))
			auth.GET("/get_all_notifications", handlers.GetAllNotifications(storage))
			auth.GET("/get_notification/:notification_id", handlers.GetNotification(storage))
			auth.DELETE("/delete_notification", handlers.DeleteNotification(storage))
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// @Summary Пинг-сервис
// @Description Проверяет доступность сервера
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "Response with pong message"
// @Router /api/ping [get]
func pingHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "pong"})
}
