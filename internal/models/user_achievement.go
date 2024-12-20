package models

import (
	"time"

	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type UserAchievements struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	AchievementID uuid.UUID
	AwardedAt     time.Time
}
