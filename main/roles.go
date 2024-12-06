package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Role struct {
	ID   int
	Name string
	Text string
}

type RoleRepository interface {
	Insert(context.Context, Role) error
	GetByID(context.Context, int) (*Role, error)
}

type PgRoleRepository struct {
	pool *pgxpool.Pool
}

func NewPgRoleRepository(pool *pgxpool.Pool) RoleRepository {
	return &PgRoleRepository{pool: pool}
}

func (r *PgRoleRepository) Insert(ctx context.Context, role Role) error {
	query := `
        INSERT INTO roles (name, text)
        VALUES ($1, $2)
    `
	_, err := r.pool.Exec(ctx, query, role.Name, role.Text)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("role with name %s already exists", role.Name)
		}
		return fmt.Errorf("error inserting role: %v", err)
	}
	return nil
}

func (r *PgRoleRepository) GetByID(ctx context.Context, id int) (*Role, error) {
	query := `
        SELECT id, name, text
        FROM roles
        WHERE id = $1
    `
	var role Role
	err := r.pool.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name, &role.Text)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role with id %d not found", id)
		}
		return nil, fmt.Errorf("error getting role: %v", err)
	}
	return &role, nil
}