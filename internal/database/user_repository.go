package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"

	"github.com/gofrs/uuid"
)

var (
	saveNewUser = `INSERT INTO users (id, name, email, password_hash, specification, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	getUserByID = `SELECT * FROM users WHERE id = $1`
)

func SaveUser(ctx context.Context, db *sql.DB, user models.User) (uuid.UUID, error) {
	_, err := db.ExecContext(ctx, saveNewUser, user.ID, user.Email, user.PasswordHash, user.Specification, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID.UUID, nil
}

func GetUserByID(ctx context.Context, db *sql.DB, id uuid.UUID) (models.User, error) {
	row := db.QueryRowContext(ctx, getUserByID, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Specification, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}
	return user, nil
}
