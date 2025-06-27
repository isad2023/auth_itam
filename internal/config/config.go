package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	UploadPath     string
	MaxFileSize    int64
	AllowedTypes   []string
}

func LoadConfig() (*AppConfig, error) {

	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}
	fmt.Println("Current working directory:", dir)

	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: could not load .env file: %v\n", err)
	}

	config := &AppConfig{
		DBUser:         getEnv("DB_USER", "itam_user"),
		DBPass:         getEnv("DB_PASSWORD", "itam_db"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBName:         getEnv("DB_NAME", "itam_auth"),
		MigrationsPath: getEnv("MIGRATIONS_PATH", ""),
		JwtSecretKey:   getEnv("JWT_SECRET_KEY", ""),
		UploadPath:     getEnv("UPLOAD_PATH", "./uploads"),
		MaxFileSize:    getEnvInt64("MAX_FILE_SIZE", 10485760), // 10MB по умолчанию
		AllowedTypes:   getEnvSlice("ALLOWED_TYPES", []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".doc", ".docx"}),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(cfg *AppConfig) error {
	missingVars := []string{}
	if cfg.DBUser == "" {
		missingVars = append(missingVars, "DB_USER")
	}
	if cfg.DBPass == "" {
		missingVars = append(missingVars, "DB_PASSWORD")
	}
	if cfg.DBName == "" {
		missingVars = append(missingVars, "DB_NAME")
	}
	if cfg.MigrationsPath == "" {
		missingVars = append(missingVars, "MIGRATIONS_PATH")
	}
	if cfg.JwtSecretKey == "" {
		missingVars = append(missingVars, "JWT_SECRET_KEY")
	}
	if cfg.UploadPath == "" {
		missingVars = append(missingVars, "UPLOAD_PATH")
	}
	if cfg.MaxFileSize == 0 {
		missingVars = append(missingVars, "MAX_FILE_SIZE")
	}
	if len(cfg.AllowedTypes) == 0 {
		missingVars = append(missingVars, "ALLOWED_TYPES")
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missingVars)
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
