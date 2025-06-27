package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"itam_auth/internal/config"
	"itam_auth/internal/models"

	"github.com/google/uuid"
)

type FileService struct {
	config *config.AppConfig
}

func NewFileService(cfg *config.AppConfig) *FileService {
	return &FileService{
		config: cfg,
	}
}

// UploadFile загружает файл и возвращает информацию о загрузке
func (fs *FileService) UploadFile(file *multipart.FileHeader, userID uuid.UUID, uploadType string, entityID *uuid.UUID) (*models.FileUpload, error) {
	// Проверяем размер файла
	if file.Size > fs.config.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", fs.config.MaxFileSize)
	}

	// Проверяем тип файла
	if !fs.isAllowedFileType(file.Filename) {
		return nil, fmt.Errorf("file type not allowed. Allowed types: %v", fs.config.AllowedTypes)
	}

	// Создаем директорию если не существует
	if err := os.MkdirAll(fs.config.UploadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(fs.config.UploadPath, fileName)

	// Открываем исходный файл
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Создаем целевой файл
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Копируем содержимое файла
	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Определяем MIME тип
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Создаем запись о загрузке
	fileUpload := &models.FileUpload{
		ID:           uuid.New(),
		UserID:       userID,
		FileName:     fileName,
		OriginalName: file.Filename,
		FilePath:     filePath,
		FileSize:     file.Size,
		MimeType:     mimeType,
		UploadType:   uploadType,
		EntityID:     entityID,
		CreatedAt:    time.Now(),
	}

	return fileUpload, nil
}

// DeleteFile удаляет файл и запись о нем
func (fs *FileService) DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetFileURL возвращает URL для доступа к файлу
func (fs *FileService) GetFileURL(fileName string) string {
	return fmt.Sprintf("/uploads/%s", fileName)
}

// isAllowedFileType проверяет, разрешен ли тип файла
func (fs *FileService) isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range fs.config.AllowedTypes {
		if strings.ToLower(allowedType) == ext {
			return true
		}
	}
	return false
}

// ValidateDocumentFile проверяет, является ли файл документом
func (fs *FileService) ValidateDocumentFile(file *multipart.FileHeader) error {
	// Проверяем MIME тип
	contentType := file.Header.Get("Content-Type")
	allowedDocTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}
	
	isAllowed := false
	for _, allowedType := range allowedDocTypes {
		if contentType == allowedType {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("document format not supported. Allowed formats: PDF, DOC, DOCX")
	}

	// Проверяем расширение файла
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedDocExts := []string{".pdf", ".doc", ".docx"}
	
	isAllowed = false
	for _, allowedExt := range allowedDocExts {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("document format not supported. Allowed formats: %v", allowedDocExts)
	}

	return nil
}

// ValidateImageFile проверяет, является ли файл изображением
func (fs *FileService) ValidateImageFile(file *multipart.FileHeader) error {
	// Проверяем MIME тип
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("file is not an image. Content-Type: %s", contentType)
	}

	// Проверяем расширение файла
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedImageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	
	isAllowed := false
	for _, allowedExt := range allowedImageExts {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("image format not supported. Allowed formats: %v", allowedImageExts)
	}

	return nil
} 