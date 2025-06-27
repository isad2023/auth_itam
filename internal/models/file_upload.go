package models

import (
	"time"

	"github.com/google/uuid"
)

type FileUpload struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	FileName     string
	OriginalName string
	FilePath     string
	FileSize     int64
	MimeType     string
	UploadType   string // 'profile_image', 'achievement_image', etc.
	EntityID     *uuid.UUID
	CreatedAt    time.Time
} 