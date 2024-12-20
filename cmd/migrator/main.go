package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	var dbUser, dbPass, dbHost, dbPort, dbName, migrationsPath, direction string

	flag.StringVar(&dbUser, "db-user", "", "database user")
	flag.StringVar(&dbPass, "db-pass", "", "database password")
	flag.StringVar(&dbHost, "db-host", "", "database host")
	flag.StringVar(&dbPort, "db-port", "", "database port")
	flag.StringVar(&dbName, "db-name", "", "database name")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&direction, "direction", "up", "migration direction: up or down")
	flag.Parse()

	loadEnv()
	if dbUser == "" {
		dbUser = os.Getenv("DB_USER")
	}
	if dbPass == "" {
		dbPass = os.Getenv("DB_PASSWORD")
	}
	if dbHost == "" {
		dbHost = os.Getenv("DB_HOST")
	}
	if dbPort == "" {
		dbPort = os.Getenv("DB_PORT")
	}
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}
	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	}

	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" || migrationsPath == "" {
		log.Fatal("Missing required arguments or environment variables")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName),
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		fmt.Println("Migrations applied")
	case "down":
		if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Failed to revert migrations: %v", err)
		}
		fmt.Println("Migrations reverted")
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'.")
	}
}
