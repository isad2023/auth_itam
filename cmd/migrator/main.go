package main

import (
	"errors"
	"flag"
	"fmt"
	"itam_auth/internal/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	flag.StringVar(&cfg.DBUser, "db-user", cfg.DBUser, "database user")
	flag.StringVar(&cfg.DBPass, "db-pass", cfg.DBPass, "database password")
	flag.StringVar(&cfg.DBHost, "db-host", cfg.DBHost, "database host")
	flag.StringVar(&cfg.DBPort, "db-port", cfg.DBPort, "database port")
	flag.StringVar(&cfg.DBName, "db-name", cfg.DBName, "database name")
	flag.StringVar(&cfg.MigrationsPath, "migrations-path", cfg.MigrationsPath, "path to migrations")
	direction := flag.String("direction", "up", "migration direction: up or down")
	flag.Parse()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		dsn,
	)

	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := applyMigration(m, *direction); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func applyMigration(m *migrate.Migrate, direction string) error {
	switch direction {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to revert migrations: %w", err)
		}
		fmt.Println("Migrations reverted successfully")
	default:
		return fmt.Errorf("invalid migration direction: %s. (Use 'up' or 'down')", direction)
	}
	return nil
}

// go run cmd/migrator/main.go --db-user="yourname" --db-pass="pass" --db-name="example" --db-host=localhost --migrations-path="./migrations"
