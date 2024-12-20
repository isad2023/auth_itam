package main

import (
	"fmt"
	"itam_auth/internal/config"
	"itam_auth/internal/database"
	"itam_auth/internal/routes"
	"log"
)

// @title LiveCode API
// @version 1.0
// @description ITaM API
// @host localhost:8080
// @BasePath /api
func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration successfully loaded.")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		appConfig.DBHost, appConfig.DBPort, appConfig.DBUser, appConfig.DBPass, appConfig.DBName)

	if err := database.Initialize(dsn); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Println("Database successfully connected.")

	router := routes.SetupRoutes()

	fmt.Println("Starting server on port 8080")
	err_server := router.Run(":8080")
	if err_server != nil {
		fmt.Println("Error starting server:", err_server)
	}

}
