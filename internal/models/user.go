package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
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
	ID            uuid.UUID
	Name          string
	Email         string
	Telegram      string
	PasswordHash  string
	PhotoURL      string
	About         string
	ResumeURL     string
	Specification Specification
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
