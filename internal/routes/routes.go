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
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://109.73.202.151:3000",
			"http://localhost:8080",
			"http://109.73.202.151:8080",
			"http://localhost:8080/auth/swagger",
			"http://109.73.202.151:8080/auth/swagger",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}
	router.Use(cors.New(config))

	auth := router.Group("/auth")
	{
		api := auth.Group("/api")
		{
			// Public routes that don't require authorization
			api.GET("/ping", pingHandler)
			api.POST("/login", handlers.Login(storage, hmacSecret))
			api.POST("/register", handlers.Register(storage))
			api.GET("/get_user/:user_id", handlers.GetUser(storage))

			// Protected routes that require authorization
			protected := api.Group("/")
			protected.Use(middleware.AuthMiddleware(hmacSecret))
			{
				protected.GET("/me", handlers.GetCurrentUser(storage))
				protected.PATCH("/update_user_info", handlers.UpdateUserInfo(storage))
				protected.GET("/get_user_roles", handlers.GetUserRoles(storage))
				protected.GET("/get_user_properties", handlers.GetUserPermissions(storage))

				//* REQUEST ROUTES
				protected.POST("/create_user_request", handlers.CreateUserRequest(storage))
				protected.GET("/get_request", handlers.GetRequest(storage))
				protected.GET("/get_all_requests", handlers.GetAllRequests(storage))
				protected.PATCH("/update_request_status", handlers.UpdateRequestStatus(storage))
				protected.DELETE("/delete_request", handlers.DeleteRequest(storage))

				//* ACHIEVEMENT ROUTES
				protected.GET("/get_user_achievements", handlers.GetAchievementsByUserID(storage))
				protected.POST("/create_achievement", handlers.CreateAchievement(storage))
				protected.PATCH("/update_achievement", handlers.UpdateAchievement(storage))
				protected.GET("/get_achievement", handlers.GetAchievementByID(storage))
				protected.GET("/get_all_achievements", handlers.GetAllAchievements(storage))
				protected.DELETE("/delete_achievement", handlers.DeleteAchievement(storage))

				//* NOTIFICATION ROUTES
				protected.POST("/create_notification", handlers.CreateNotification(storage))
				protected.PATCH("/update_notification", handlers.UpdateNotification(storage))
				protected.GET("/get_all_notifications", handlers.GetAllNotifications(storage))
				protected.GET("/get_notification/:notification_id", handlers.GetNotification(storage))
				protected.DELETE("/delete_notification", handlers.DeleteNotification(storage))
			}
		}

		auth.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}

// @Summary Пинг-сервис
// @Description Проверяет доступность сервера
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "Response with pong message"
// @Router /auth/api/ping [get]
func pingHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "pong"})
}
