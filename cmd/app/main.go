package main

import (
	"fmt"
	_ "itam_auth/docs"
	"itam_auth/internal/config"
	"itam_auth/internal/database"
	"itam_auth/internal/routes"
	"log"
)

// @title LiveCode API
// @version 1.0
// @description ITaM API
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for authorization. Format: "Bearer <token>"
// @BasePath /api
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

	router := routes.SetupRoutes(storage, appConfig.JwtSecretKey)
	log.Printf("Starting server on port %s", serverPort)
	if err := router.Run(serverPort); err != nil {
		fmt.Printf("Error starting server: %v", err)
	}

}
