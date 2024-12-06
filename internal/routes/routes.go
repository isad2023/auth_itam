package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	config.AllowMethods = []string{"*"}
	router.Use(cors.New(config))

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//* USER
	router.POST("/api/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "login",
		})
	})

	router.GET("/api/get_user", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_user",
		})
	})

	router.PATCH("/api/update_user_info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "update_user_info",
		})
	})

	router.GET("/api/get_user_roles", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_user_roles",
		})
	})

	router.GET("/api/get_user_properties", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_user_properties",
		})
	})

	router.GET("/api/get_user_achievements", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_user_achievements",
		})
	})

	router.GET("/api/get_all_user_notifications", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_all_user_notifications",
		})
	})

	//* Requests

	router.POST("/api/create_user_request", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "create_user_request",
		})
	})

	router.GET("/api/get_request", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_request",
		})
	})

	router.GET("/api/get_all_requests", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_all_requests",
		})
	})

	router.PATCH("/api/update_request_status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "update_request_status",
		})
	})

	//* Achievements

	router.POST("/api/create_achievement", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "create_achievement",
		})
	})

	router.PATCH("/api/update_achievement", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "update_achievement",
		})
	})

	router.GET("/api/get_achievement", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_achievement",
		})
	})

	router.GET("/api/get_all_achievements", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_all_achievements",
		})
	})

	//* Notifications (admin)

	router.POST("/api/create_notification", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "create_notification",
		})
	})

	router.PATCH("/api/update_notification", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "update_notification",
		})
	})

	router.GET("/api/get_all_notifications", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_all_notifications",
		})
	})

	router.GET("/api/get_notification", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_notification",
		})
	})

	return router
}
