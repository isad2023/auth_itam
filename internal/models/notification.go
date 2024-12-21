package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Content   string
	IsRead    bool
	CreatedAt time.Time
}
