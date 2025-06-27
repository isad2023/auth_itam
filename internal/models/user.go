package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Specification string

const (
	Frontend        Specification = "Frontend"
	Backend         Specification = "Backend"
	MachineLearning Specification = "Machine Learning"
	Developer       Specification = "Developer"
	Designer        Specification = "Designer"
	Manager         Specification = "Manager"
)

func (s *Specification) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("invalid data type for Specification")
	}

	switch strValue {
	case string(Frontend), string(Backend), string(MachineLearning), string(Developer), string(Designer), string(Manager):
		*s = Specification(strValue)
		return nil
	default:
		return fmt.Errorf("invalid value '%s' for specification", strValue)
	}
}

func (s Specification) Value() (driver.Value, error) {
	switch s {
	case Frontend, Backend, MachineLearning, Developer, Designer, Manager:
		return string(s), nil
	default:
		return nil, fmt.Errorf("invalid value '%s' for specification", string(s))
	}
}

type User struct {
	ID            uuid.UUID      `json:"ID"`
	Name          string         `json:"Name"`
	Email         string         `json:"Email"`
	Telegram      *string        `json:"Telegram"`
	PasswordHash  string         `json:"PasswordHash"`
	PhotoURL      *string        `json:"PhotoURL"`
	About         *string        `json:"About"`
	ResumeURL     *string        `json:"ResumeURL"`
	Specification Specification `json:"Specification"`
	CreatedAt     time.Time      `json:"CreatedAt"`
	UpdatedAt     time.Time      `json:"UpdatedAt"`
}

func (u *User) GetAdminServices(userRoles []UserRole, roles []Role, rolePermissions []RolePermission, permissions []Permission) []string {
	adminServices := make(map[string]struct{})

	roleMap := make(map[uuid.UUID]Role)
	for _, r := range roles {
		roleMap[r.ID] = r
	}

	rolePermMap := make(map[uuid.UUID][]uuid.UUID)
	for _, rp := range rolePermissions {
		rolePermMap[rp.RoleID] = append(rolePermMap[rp.RoleID], rp.PermissionID)
	}

	permMap := make(map[uuid.UUID]Permission)
	for _, p := range permissions {
		permMap[p.ID] = p
	}

	for _, ur := range userRoles {
		if ur.UserID != u.ID {
			continue
		}
		for _, permID := range rolePermMap[ur.RoleID] {
			if perm, ok := permMap[permID]; ok && strings.HasPrefix(perm.Name, "admin_") {
				service := strings.TrimPrefix(perm.Name, "admin_")
				adminServices[service] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(adminServices))
	for service := range adminServices {
		result = append(result, service)
	}
	return result
}
