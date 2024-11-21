package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
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
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserRepository interface {
	Insert(context.Context, User) error
	GetByID(context.Context, int) (*User, error)
}

type PgUserRepository struct {
	pool *pgxpool.Pool
}

func NewPgUserRepository(pool *pgxpool.Pool) UserRepository {
	return &PgUserRepository{pool: pool}
}

func (r *PgUserRepository) Insert(ctx context.Context, user User) error {
	query := `
        INSERT INTO users (name, email, telegram, password_hash, photo_url, about, resume_url, specification)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := r.pool.Exec(ctx, query,
		user.Name,
		user.Email,
		user.Telegram,
		user.PasswordHash,
		user.PhotoURL,
		user.About,
		user.ResumeURL,
		user.Specification,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("email already exists: %v", user.Email)
		}
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func (r *PgUserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `
        SELECT id, name, email, telegram, password_hash, photo_url, about, resume_url, specification, created_at, updated_at
        FROM users
        WHERE id = $1
    `
	var user User
	err := r.pool.QueryRow(ctx, query, id).Scan(
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	pool, err := pgxpool.Connect(context.Background(), psqlInfo)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	repo := NewPgUserRepository(pool)

	err = repo.Insert(context.Background(), User{
	    Name:          "John Doe",
	    Email:         "john.doe@example.com",
	    Telegram:      "john_telegram",
	    PasswordHash:  "password_hash",
	    PhotoURL:      "http://example.com/photo.jpg",
	    About:         "About John",
	    ResumeURL:     "http://example.com/resume.pdf",
	    Specification: "Developer",
	})
	if err != nil {
	    log.Fatalf("Error inserting user: %v", err)
	}
	fmt.Println("User successfully inserted!")

	user, err := repo.GetByID(context.Background(), 1)
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
	}
	fmt.Printf("User: %+v\n", user)
}
