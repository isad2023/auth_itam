package database

import (
	"context"
	"itam_auth/internal/models"

	"github.com/google/uuid"
)

// SaveFileUpload сохраняет информацию о загруженном файле
func (s *Storage) SaveFileUpload(ctx context.Context, fileUpload *models.FileUpload) error {
	query := `
		INSERT INTO file_uploads (id, user_id, file_name, original_name, file_path, file_size, mime_type, upload_type, entity_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := s.db.ExecContext(ctx, query,
		fileUpload.ID,
		fileUpload.UserID,
		fileUpload.FileName,
		fileUpload.OriginalName,
		fileUpload.FilePath,
		fileUpload.FileSize,
		fileUpload.MimeType,
		fileUpload.UploadType,
		fileUpload.EntityID,
		fileUpload.CreatedAt,
	)

	return err
}

// GetFileUploadByID получает информацию о файле по ID
func (s *Storage) GetFileUploadByID(ctx context.Context, id uuid.UUID) (*models.FileUpload, error) {
	query := `
		SELECT id, user_id, file_name, original_name, file_path, file_size, mime_type, upload_type, entity_id, created_at
		FROM file_uploads
		WHERE id = $1
	`

	var fileUpload models.FileUpload
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&fileUpload.ID,
		&fileUpload.UserID,
		&fileUpload.FileName,
		&fileUpload.OriginalName,
		&fileUpload.FilePath,
		&fileUpload.FileSize,
		&fileUpload.MimeType,
		&fileUpload.UploadType,
		&fileUpload.EntityID,
		&fileUpload.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &fileUpload, nil
}

// GetFileUploadsByUserID получает все файлы пользователя
func (s *Storage) GetFileUploadsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.FileUpload, error) {
	query := `
		SELECT id, user_id, file_name, original_name, file_path, file_size, mime_type, upload_type, entity_id, created_at
		FROM file_uploads
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fileUploads []*models.FileUpload
	for rows.Next() {
		var fileUpload models.FileUpload
		err := rows.Scan(
			&fileUpload.ID,
			&fileUpload.UserID,
			&fileUpload.FileName,
			&fileUpload.OriginalName,
			&fileUpload.FilePath,
			&fileUpload.FileSize,
			&fileUpload.MimeType,
			&fileUpload.UploadType,
			&fileUpload.EntityID,
			&fileUpload.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		fileUploads = append(fileUploads, &fileUpload)
	}

	return fileUploads, nil
}

// GetFileUploadsByTypeAndEntity получает файлы по типу и ID сущности
func (s *Storage) GetFileUploadsByTypeAndEntity(ctx context.Context, uploadType string, entityID uuid.UUID) ([]*models.FileUpload, error) {
	query := `
		SELECT id, user_id, file_name, original_name, file_path, file_size, mime_type, upload_type, entity_id, created_at
		FROM file_uploads
		WHERE upload_type = $1 AND entity_id = $2
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, uploadType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fileUploads []*models.FileUpload
	for rows.Next() {
		var fileUpload models.FileUpload
		err := rows.Scan(
			&fileUpload.ID,
			&fileUpload.UserID,
			&fileUpload.FileName,
			&fileUpload.OriginalName,
			&fileUpload.FilePath,
			&fileUpload.FileSize,
			&fileUpload.MimeType,
			&fileUpload.UploadType,
			&fileUpload.EntityID,
			&fileUpload.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		fileUploads = append(fileUploads, &fileUpload)
	}

	return fileUploads, nil
}

// DeleteFileUpload удаляет запись о файле
func (s *Storage) DeleteFileUpload(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM file_uploads WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// UpdateUserProfileImage обновляет URL изображения профиля пользователя
func (s *Storage) UpdateUserProfileImage(ctx context.Context, userID uuid.UUID, imageURL string) error {
	query := `UPDATE users SET photo_url = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, imageURL, userID)
	return err
}

// UpdateAchievementImage обновляет URL изображения достижения
func (s *Storage) UpdateAchievementImage(ctx context.Context, achievementID uuid.UUID, imageURL string) error {
	query := `UPDATE achievements SET image_url = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, imageURL, achievementID)
	return err
}

// UpdateUserResumeURL обновляет URL резюме пользователя
func (s *Storage) UpdateUserResumeURL(ctx context.Context, userID uuid.UUID, resumeURL string) error {
	query := `UPDATE users SET resume_url = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, resumeURL, userID)
	return err
} 