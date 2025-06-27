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

// User представляет пользователя системы
type User struct {
	ID            uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string    `json:"name" example:"John Doe"`
	Email         string    `json:"email" example:"john@example.com"`
	Telegram      *string   `json:"telegram,omitempty" example:"@johndoe"`
	PasswordHash  string    `json:"-"` // Не отображается в JSON
	PhotoURL      *string   `json:"photo_url,omitempty" example:"/uploads/profile.jpg"`
	About         *string   `json:"about,omitempty" example:"Software developer with 5 years of experience"`
	ResumeURL     *string   `json:"resume_url,omitempty" example:"/uploads/resume.pdf"`
	Specification Specification `json:"specification" example:"Backend"`
	CreatedAt     time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
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
