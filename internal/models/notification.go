package models

import (
	"time"

	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type Notification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Content   string
	IsRead    bool
	CreatedAt time.Time
}
