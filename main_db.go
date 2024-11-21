package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "itam_user"
	dbname   = "itam_auth"
	password = "itam_db"
)

type User struct {
	ID            int
	Name          string
	Email         string
	Telegram      string
	PasswordHash  string
	PhotoURL      string
	About         string
	ResumeURL     string
	Specification string
	CreatedAt     string
	UpdatedAt     string
}

func InsertUser(db *sql.DB, name, email, telegram, password, photoURL, about, resumeURL, specification string) error {
	query := `
        INSERT INTO users (name, email, telegram, password_hash, photo_url, about, resume_url, specification)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := db.Exec(query, name, email, telegram, password, photoURL, about, resumeURL, specification)

	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `
        SELECT id, name, email, telegram, password_hash, photo_url, about, resume_url, specification, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user User
	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Telegram,
		&user.PasswordHash,
		&user.PhotoURL,
		&user.About,
		&user.ResumeURL,
		&user.Specification,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &user, nil
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	err = InsertUser(db, "John Doe", "john.doe@example.com", "john_telegram", "password_hash", "http://example.com/photo.jpg", "About John", "http://example.com/resume.pdf", "Developer")
	if err != nil {
		log.Fatalf("Error inserting user: %v", err)
	}

	fmt.Println("User successfully inserted!")

}
