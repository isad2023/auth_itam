package main

import (
	"fmt"
	_ "itam_auth/docs"
	"itam_auth/internal/config"
	"itam_auth/internal/database"
	"itam_auth/internal/routes"
	"log"
)

// @title ITaM Auth API
// @version 1.0
// @description API для системы аутентификации и управления пользователями ITaM
// @description Включает функциональность для работы с пользователями, достижениями, запросами, уведомлениями и загрузкой файлов
// @host 109.73.202.151:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl /auth/api/login
// @in header
// @name Authorization
// @contact.name ITaM Support
// @contact.email support@itam.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
const (
	serverPort = ":8080"
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration successfully loaded.")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBUser,
		appConfig.DBPass,
		appConfig.DBName,
	)

	storage, err := database.Initialize(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer storage.Close()
	log.Println("Database successfully connected.")

	router := routes.SetupRoutes(storage, appConfig.JwtSecretKey, appConfig)
	log.Printf("Starting server on port %s", serverPort)
	if err := router.Run(serverPort); err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}
