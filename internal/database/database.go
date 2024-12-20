package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func Initialize(dsn string) (*Storage, error) {
	var err error

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to database!")
	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
}
