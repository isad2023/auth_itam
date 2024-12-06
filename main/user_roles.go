package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRole struct {
	UserID int
	RoleID int
}

type UserRoleRepository interface {
	Insert(context.Context, UserRole) error
	GetByUserID(context.Context, int) ([]UserRole, error)
}

type PgUserRoleRepository struct {
	pool *pgxpool.Pool
}

func NewPgUserRoleRepository(pool *pgxpool.Pool) UserRoleRepository {
	return &PgUserRoleRepository{pool: pool}
}

func (r *PgUserRoleRepository) Insert(ctx context.Context, userRole UserRole) error {
	query := `
        INSERT INTO user_roles (user_id, role_id)
        VALUES ($1, $2)
    `
	_, err := r.pool.Exec(ctx, query, userRole.UserID, userRole.RoleID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("user role already exists: user_id=%d, role_id=%d", userRole.UserID, userRole.RoleID)
		}
		return fmt.Errorf("error inserting user role: %v", err)
	}
	return nil
}

func (r *PgUserRoleRepository) GetByUserID(ctx context.Context, userID int) ([]UserRole, error) {
	query := `
        SELECT user_id, role_id
        FROM user_roles
        WHERE user_id = $1
    `
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user roles: %v", err)
	}
	defer rows.Close()

	var userRoles []UserRole
	for rows.Next() {
		var userRole UserRole
		err := rows.Scan(&userRole.UserID, &userRole.RoleID)
		if err != nil {
			return nil, fmt.Errorf("error scanning user role: %v", err)
		}
		userRoles = append(userRoles, userRole)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user roles: %v", err)
	}

	return userRoles, nil
}
