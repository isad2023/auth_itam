package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"

	"github.com/google/uuid"
)

var (
	saveNewRole            = `INSERT INTO roles (id, name) VALUES ($1, $2)`
	getRoleByID            = `SELECT * FROM roles WHERE id = $1`
	saveNewPermission      = `INSERT INTO permissions (id, name) VALUES ($1, $2)`
	getPermissionByID      = `SELECT * FROM permissions WHERE id = $1`
	getPermissionsByRoleID = `SELECT p.id, p.name FROM permissions p INNER JOIN role_permissions rp ON p.id = rp.permission_id WHERE rp.role_id = $1`
	getRolesByUserID       = `SELECT r.id, r.name FROM roles r INNER JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1`
	getUserPermissions     = `
		SELECT p.id, p.name
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1`
)

func (s *Storage) SaveRole(ctx context.Context, role models.Role) (uuid.UUID, error) {
	_, err := s.db.ExecContext(ctx, saveNewRole, role.ID, role.Name)
	if err != nil {
		return uuid.Nil, err
	}
	return role.ID, nil
}

func (s *Storage) GetRole(ctx context.Context, id uuid.UUID) (models.Role, error) {
	row := s.db.QueryRowContext(ctx, getRoleByID, id)

	var role models.Role
	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return role, fmt.Errorf("role not found")
		}
		return role, err
	}
	return role, nil
}

func (s *Storage) SavePermission(ctx context.Context, permission models.Permission) (uuid.UUID, error) {
	_, err := s.db.ExecContext(ctx, saveNewPermission, permission.ID, permission.Name)
	if err != nil {
		return uuid.Nil, err
	}
	return permission.ID, nil
}

func (s *Storage) GetPermission(ctx context.Context, id uuid.UUID) (models.Permission, error) {
	row := s.db.QueryRowContext(ctx, getPermissionByID, id)

	var permission models.Permission
	err := row.Scan(&permission.ID, &permission.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return permission, fmt.Errorf("permission not found")
		}
		return permission, err
	}
	return permission, nil
}

func (s *Storage) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.RolePermission, error) {
	rows, err := s.db.QueryContext(ctx, getPermissionsByRoleID, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rolePermissions []models.RolePermission
	for rows.Next() {
		var rolePermission models.RolePermission
		err := rows.Scan(&rolePermission.ID, &rolePermission.RoleID, &rolePermission.PermissionID)
		if err != nil {
			return nil, err
		}
		rolePermissions = append(rolePermissions, rolePermission)
	}
	return rolePermissions, nil
}

func (s *Storage) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.UserRole, error) {
	rows, err := s.db.QueryContext(ctx, getRolesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRoles []models.UserRole
	for rows.Next() {
		var userRole models.UserRole
		err := rows.Scan(&userRole.ID, &userRole.UserID, &userRole.RoleID)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, userRole)
	}
	return userRoles, nil
}

func (s *Storage) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]models.Permission, error) {
	rows, err := s.db.QueryContext(ctx, getUserPermissions, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		err := rows.Scan(&permission.ID, &permission.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return permissions, nil
}
