package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

var (
	saveNewRole            = `INSERT INTO roles (id, name) VALUES ($1, $2)`
	getRoleByID            = `SELECT id, name FROM roles WHERE id = $1`
	getRoleByName          = `SELECT id, name FROM roles WHERE name = $1`
	saveNewPermission      = `INSERT INTO permissions (id, name) VALUES ($1, $2)`
	getPermissionByID      = `SELECT id, name FROM permissions WHERE id = $1`
	getPermissionsByRoleID = `SELECT p.id, p.name FROM permissions p INNER JOIN role_permissions rp ON p.id = rp.permission_id WHERE rp.role_id = $1`
	getRolesByUserID       = `SELECT r.id, r.name FROM roles r INNER JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1`
	getUserPermissions     = `
		SELECT p.id, p.name
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1`
	saveUserRole        = `INSERT INTO user_roles (id, user_id, role_id) VALUES ($1, $2, $3)`
	getRolesByIDs       = `SELECT id, name FROM roles WHERE id = ANY($1)`
	getPermissionsByIDs = `SELECT id, name FROM permissions WHERE id = ANY($1)`
	getRolePermissions  = `SELECT id, role_id, permission_id FROM role_permissions WHERE role_id = $1`
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

func (s *Storage) GetRoleByName(ctx context.Context, name string) (models.Role, error) {
	row := s.db.QueryRowContext(ctx, getRoleByName, name)

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
	rows, err := s.db.QueryContext(ctx, getRolePermissions, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	defer rows.Close()

	var rolePermissions []models.RolePermission
	for rows.Next() {
		var rp models.RolePermission
		err := rows.Scan(&rp.ID, &rp.RoleID, &rp.PermissionID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role permission: %w", err)
		}
		rolePermissions = append(rolePermissions, rp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return rolePermissions, nil
}

func (s *Storage) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.UserRole, error) {
	rows, err := s.db.QueryContext(ctx, getRolesByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var userRoles []models.UserRole
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user role: %w", err)
		}
		userRoles = append(userRoles, models.UserRole{
			ID:     uuid.New(),
			UserID: userID,
			RoleID: role.ID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
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

func (s *Storage) SaveUserRole(ctx context.Context, userRole models.UserRole) (uuid.UUID, error) {
	_, err := s.db.ExecContext(ctx, saveUserRole, userRole.ID, userRole.UserID, userRole.RoleID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to save user role: %w", err)
	}
	return userRole.ID, nil
}

func (s *Storage) GetRolesByIDs(ctx context.Context, roleIDs []uuid.UUID) ([]models.Role, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}

	rows, err := s.db.QueryContext(ctx, getRolesByIDs, pq.Array(roleIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get roles by IDs: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return roles, nil
}

func (s *Storage) GetPermissionsByIDs(ctx context.Context, permissionIDs []uuid.UUID) ([]models.Permission, error) {
	if len(permissionIDs) == 0 {
		return nil, nil
	}

	rows, err := s.db.QueryContext(ctx, getPermissionsByIDs, pq.Array(permissionIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by IDs: %w", err)
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
