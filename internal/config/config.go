package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBUser         string
	DBPass         string
	DBHost         string
	DBPort         string
	DBName         string
	MigrationsPath string
	JwtSecretKey   string
}

func LoadConfig() (*AppConfig, error) {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	fmt.Println("Current working directory:", dir)

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	config := &AppConfig{
		DBUser:         os.Getenv("DB_USER"),
		DBPass:         os.Getenv("DB_PASSWORD"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBName:         os.Getenv("DB_NAME"),
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		JwtSecretKey:   os.Getenv("JWT_SECRET_KEY"),
	}

	// Проверяем обязательные параметры
	missingVars := []string{}
	if config.DBUser == "" {
		missingVars = append(missingVars, "DB_USER")
	}
	if config.DBPass == "" {
		missingVars = append(missingVars, "DB_PASSWORD")
	}
	if config.DBName == "" {
		missingVars = append(missingVars, "DB_NAME")
	}
	if config.MigrationsPath == "" {
		missingVars = append(missingVars, "MIGRATIONS_PATH")
	}

	if config.JwtSecretKey == "" {
		missingVars = append(missingVars, "JWT_SECRET_KEY")
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missingVars)
	}

	return config, nil
}
