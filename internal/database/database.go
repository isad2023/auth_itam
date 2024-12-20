package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Initialize(dsn string) error {
	var err error
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to database!")
	return nil
}

func Close() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
}
