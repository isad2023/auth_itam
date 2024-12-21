package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"

	"github.com/google/uuid"
)

const (
	saveNewUserQuery = `INSERT INTO users (id, name, email, password_hash, specification, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	getUserByIDQuery    = `SELECT id, name, email, password_hash, specification, created_at, updated_at FROM users WHERE id = $1`
	getUserByEmailQuery = `SELECT id, name, email, password_hash FROM users WHERE email = $1`
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
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Specification, &user.CreatedAt, &user.UpdatedAt)
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
