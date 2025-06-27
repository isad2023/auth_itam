package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"
	"log"

	"github.com/google/uuid"
)

const (
	saveNewUserQuery = `INSERT INTO users (id, name, email, password_hash, specification, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	getUserByIDQuery    = `SELECT id, name, email, password_hash, specification, about, photo_url, resume_url, telegram, created_at, updated_at FROM users WHERE id = $1`
	getUserByEmailQuery = `SELECT id, name, email, password_hash FROM users WHERE email = $1`
	updateUserQuery     = `UPDATE users SET name = $1, specification = $2, about = $3, photo_url = $4, resume_url = $5, telegram = $6, updated_at = $7 WHERE id = $8`
)

func (s *Storage) SaveUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	_, err := s.db.ExecContext(ctx, saveNewUserQuery,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Specification,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	row := s.db.QueryRowContext(ctx, getUserByIDQuery, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Specification,
		&user.About,
		&user.PhotoURL,
		&user.ResumeURL,
		&user.Telegram,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}
	return user, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	row := s.db.QueryRowContext(ctx, getUserByEmailQuery, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}
	return user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user models.User) error {
	if user.ID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction for updating user with ID %s: %v", user.ID, err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Failed to rollback transaction for user with ID %s: %v", user.ID, err)
		}
	}()

	result, err := tx.ExecContext(ctx, updateUserQuery,
		user.Name,
		user.Specification,
		user.About,
		user.PhotoURL,
		user.ResumeURL,
		user.Telegram,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		log.Printf("Failed to update user with ID %s: %v", user.ID, err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected for user with ID %s: %v", user.ID, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID: %s", user.ID)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction for user with ID %s: %v", user.ID, err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
