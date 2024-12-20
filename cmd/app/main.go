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
// @BasePath /api
func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration successfully loaded.")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		appConfig.DBHost, appConfig.DBPort, appConfig.DBUser, appConfig.DBPass, appConfig.DBName)

	storage, err := database.Initialize(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer storage.Close()
	log.Println("Database successfully connected.")

	router := routes.SetupRoutes(storage)

	fmt.Println("Starting server on port 8080")
	err_server := router.Run(":8080")
	if err_server != nil {
		fmt.Println("Error starting server:", err_server)
	}

}
