package database

import (
	"context"
	"database/sql"
	"fmt"
	"itam_auth/internal/models"

	"github.com/gofrs/uuid"
)

var (
	saveNewRole = `INSERT INTO roles (id, name) VALUES ($1, $2)`
	getRoleByID = `SELECT * FROM roles WHERE id = $1`

	saveNewPermission = `INSERT INTO permissions (id, name) VALUES ($1, $2)`
	getPermissionByID = `SELECT * FROM permissions WHERE id = $1`

	saveNewRolePermission  = `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)`
	getPermissionsByRoleID = `SELECT p.id, p.name FROM permissions p INNER JOIN role_permissions rp ON p.id = rp.permission_id WHERE rp.role_id = $1`

	saveNewUserRole  = `INSERT INTO user_roles (id, user_id, role_id) VALUES ($1, $2, $3)`
	getRolesByUserID = `SELECT r.id, r.name FROM roles r INNER JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1`
)

func SaveRole(ctx context.Context, db *sql.DB, role models.Role) (uuid.UUID, error) {
	_, err := db.ExecContext(ctx, saveNewRole, role.ID, role.Name)
	if err != nil {
		return uuid.Nil, err
	}
	return role.ID.UUID, nil
}

func GetRole(ctx context.Context, db *sql.DB, id uuid.UUID) (models.Role, error) {
	row := db.QueryRowContext(ctx, getRoleByID, id)

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

func SavePermission(ctx context.Context, db *sql.DB, permission models.Permission) (uuid.UUID, error) {
	_, err := db.ExecContext(ctx, saveNewPermission, permission.ID, permission.Name)
	if err != nil {
		return uuid.Nil, err
	}
	return permission.ID.UUID, nil
}

func GetPermission(ctx context.Context, db *sql.DB, id uuid.UUID) (models.Permission, error) {
	row := db.QueryRowContext(ctx, getPermissionByID, id)

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


func GetRolePermissions(ctx context.Context, db *sql.DB, roleID uuid.UUID) ([]models.RolePermission, error) {
	rows, err := db.QueryContext(ctx, getPermissionsByRoleID, roleID)
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



func GetUserRoles(ctx context.Context, db *sql.DB, userID uuid.UUID) ([]models.UserRole, error) {
	rows, err := db.QueryContext(ctx, getRolesByUserID, userID)
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
