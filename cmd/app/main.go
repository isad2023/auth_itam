package main

import (
	"fmt"
	"itam_auth/internal/routes"
)

// @title LiveCode API
// @version 1.0
// @description ITaM API
// @host localhost:8080
// @BasePath /api
func main() {
	router := routes.SetupRoutes()

	fmt.Println("Starting server on port 8080")
	err_server := router.Run(":8080")
	if err_server != nil {
		fmt.Println("Error starting server:", err_server)
	}

}
