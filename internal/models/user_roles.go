package models

import uuid "github.com/jackc/pgtype/ext/gofrs-uuid"

type UserRole struct {
	ID     uuid.UUID
	UserID uuid.UUID
	RoleID uuid.UUID
}

type Role struct {
	ID   uuid.UUID
	Name string
}

type RolePermission struct {
	ID           uuid.UUID
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}

type Permission struct {
	ID   uuid.UUID
	Name string
}
