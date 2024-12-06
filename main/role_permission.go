package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RolePermission struct {
	RoleID       int
	PermissionID int
}

type RolePermissionRepository interface {
	Insert(context.Context, RolePermission) error
	GetByRoleID(context.Context, int) ([]RolePermission, error)
}

type PgRolePermissionRepository struct {
	pool *pgxpool.Pool
}

func NewPgRolePermissionRepository(pool *pgxpool.Pool) RolePermissionRepository {
	return &PgRolePermissionRepository{pool: pool}
}

func (r *PgRolePermissionRepository) Insert(ctx context.Context, rolePermission RolePermission) error {
	query := `
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES ($1, $2)
    `
	_, err := r.pool.Exec(ctx, query, rolePermission.RoleID, rolePermission.PermissionID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("role permission already exists: role_id=%d, permission_id=%d", rolePermission.RoleID, rolePermission.PermissionID)
		}
		return fmt.Errorf("error inserting role permission: %v", err)
	}
	return nil
}

func (r *PgRolePermissionRepository) GetByRoleID(ctx context.Context, roleID int) ([]RolePermission, error) {
	query := `
        SELECT role_id, permission_id
        FROM role_permissions
        WHERE role_id = $1
    `
	rows, err := r.pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("error querying role permissions: %v", err)
	}
	defer rows.Close()

	var rolePermissions []RolePermission
	for rows.Next() {
		var rolePermission RolePermission
		err := rows.Scan(&rolePermission.RoleID, &rolePermission.PermissionID)
		if err != nil {
			return nil, fmt.Errorf("error scanning role permission: %v", err)
		}
		rolePermissions = append(rolePermissions, rolePermission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over role permissions: %v", err)
	}

	return rolePermissions, nil
}
