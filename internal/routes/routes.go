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

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	config.AllowMethods = []string{"*"}
	router.Use(cors.New(config))

	api := router.Group("/api")
	{
		api.GET("/ping", pingHandler)

		//* USER ROUTES
		api.POST("/login", handlers.Login(storage, hmacSecret))
		api.POST("/register", handlers.Register(storage))
		api.GET("/get_user/:user_id", handlers.GetUser(storage))

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware(hmacSecret))
		{
			api.PATCH("/update_user_info", updateUserInfoHandler())
			api.GET("/get_user_roles", handlers.GetUserRoles(storage))
			api.GET("/get_user_properties", handlers.GetUserPermissions(storage))

			//* REQUEST ROUTES
			api.POST("/create_user_request", handlers.CreateUserRequest(storage))
			api.GET("/get_request", handlers.GetRequest(storage))
			api.GET("/get_all_requests", handlers.GetAllRequests(storage))
			api.PATCH("/update_request_status", handlers.UpdateRequestStatus(storage))

			//* ACHIEVEMENT ROUTES
			api.GET("/get_user_achievements", handlers.GetAchievementsByUserID(storage))
			api.POST("/create_achievement", handlers.CreateAchievement(storage))
			api.PATCH("/update_achievement", handlers.UpdateAchievement(storage))
			api.GET("/get_achievement", handlers.GetAchievementByID(storage))
			api.GET("/get_all_achievements", handlers.GetAllAchievements(storage))
			api.DELETE("/delete_achievement", handlers.DeleteAchievement(storage))

			//* NOTIFICATION ROUTES
			api.POST("/create_notification", handlers.CreateNotification(storage))
			api.PATCH("/update_notification", handlers.UpdateNotification(storage))
			api.GET("/get_all_notifications", handlers.GetAllNotifications(storage))
			api.GET("/get_notification/:notification_id", handlers.GetNotification(storage))
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

// @Summary Обновить информацию пользователя
// @Description Обновляет профиль пользователя
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Success message"
// @Router /api/update_user_info [patch]
func updateUserInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "update_user_info"})
	}
}
