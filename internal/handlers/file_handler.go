package handlers

import (
	"itam_auth/internal/config"
	"itam_auth/internal/database"
	"itam_auth/internal/services/file"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Загрузить изображение профиля
// @Description Загружает изображение профиля для текущего пользователя и обновляет поле photo_url
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Profile image file (JPEG, PNG, GIF, WebP, max 10MB)"
// @Security OAuth2Password
// @Success 200 {object} models.FileUploadResponse "Success message with file info"
// @Failure 400 {object} models.ErrorResponse "Invalid file or request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/api/upload_profile_image [post]
func UploadProfileImage(storage *database.Storage, fileService *file.FileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID пользователя из контекста
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Преобразуем строку в UUID
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Получаем файл из запроса
		fileHeader, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
			return
		}

		// Проверяем, что файл является изображением
		if err := fileService.ValidateImageFile(fileHeader); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Загружаем файл
		fileUpload, err := fileService.UploadFile(fileHeader, userID, "profile_image", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}

		// Сохраняем информацию о файле в базе данных
		ctx := c.Request.Context()
		if err := storage.SaveFileUpload(ctx, fileUpload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file information"})
			return
		}

		// Обновляем URL изображения профиля пользователя
		imageURL := fileService.GetFileURL(fileUpload.FileName)
		if err := storage.UpdateUserProfileImage(ctx, userID, imageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile image URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Profile image uploaded successfully",
			"file": gin.H{
				"id":           fileUpload.ID,
				"fileName":     fileUpload.FileName,
				"originalName": fileUpload.OriginalName,
				"fileSize":     fileUpload.FileSize,
				"mimeType":     fileUpload.MimeType,
				"url":          imageURL,
			},
		})
	}
}

// @Summary Загрузить изображение достижения
// @Description Загружает изображение для достижения и обновляет поле image_url
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Achievement image file (JPEG, PNG, GIF, WebP, max 10MB)"
// @Param achievement_id formData string true "Achievement ID (UUID)"
// @Security OAuth2Password
// @Success 200 {object} models.FileUploadResponse "Success message with file info"
// @Failure 400 {object} models.ErrorResponse "Invalid file or request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/api/upload_achievement_image [post]
func UploadAchievementImage(storage *database.Storage, fileService *file.FileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID пользователя из контекста
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Преобразуем строку в UUID
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Получаем файл из запроса
		fileHeader, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
			return
		}

		// Получаем ID достижения
		achievementIDStr := c.PostForm("achievement_id")
		if achievementIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Achievement ID is required"})
			return
		}

		achievementID, err := uuid.Parse(achievementIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid achievement ID"})
			return
		}

		// Проверяем, что файл является изображением
		if err := fileService.ValidateImageFile(fileHeader); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Загружаем файл
		fileUpload, err := fileService.UploadFile(fileHeader, userID, "achievement_image", &achievementID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}

		// Сохраняем информацию о файле в базе данных
		ctx := c.Request.Context()
		if err := storage.SaveFileUpload(ctx, fileUpload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file information"})
			return
		}

		// Обновляем URL изображения достижения
		imageURL := fileService.GetFileURL(fileUpload.FileName)
		if err := storage.UpdateAchievementImage(ctx, achievementID, imageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievement image URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Achievement image uploaded successfully",
			"file": gin.H{
				"id":           fileUpload.ID,
				"fileName":     fileUpload.FileName,
				"originalName": fileUpload.OriginalName,
				"fileSize":     fileUpload.FileSize,
				"mimeType":     fileUpload.MimeType,
				"url":          imageURL,
			},
		})
	}
}

// @Summary Получить файл
// @Description Возвращает загруженный файл по имени
// @Tags Files
// @Produce octet-stream
// @Param filename path string true "File name (UUID format)"
// @Success 200 {file} file "File content"
// @Failure 400 {object} map[string]string "Invalid filename"
// @Failure 404 {object} map[string]string "File not found"
// @Router /uploads/{filename} [get]
func ServeFile(cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
			return
		}

		// Проверяем безопасность пути
		if filepath.Ext(filename) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
			return
		}

		filePath := filepath.Join(cfg.UploadPath, filename)
		c.File(filePath)
	}
}

// @Summary Получить список файлов пользователя
// @Description Возвращает список всех файлов, загруженных пользователем
// @Tags Files
// @Produce json
// @Security OAuth2Password
// @Success 200 {array} models.FileInfo "List of user files"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/api/get_user_files [get]
func GetUserFiles(storage *database.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID пользователя из контекста
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Преобразуем строку в UUID
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Получаем список файлов пользователя
		ctx := c.Request.Context()
		fileUploads, err := storage.GetFileUploadsByUserID(ctx, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user files"})
			return
		}

		// Преобразуем в формат ответа
		var files []gin.H
		for _, fileUpload := range fileUploads {
			files = append(files, gin.H{
				"id":           fileUpload.ID,
				"fileName":     fileUpload.FileName,
				"originalName": fileUpload.OriginalName,
				"fileSize":     fileUpload.FileSize,
				"mimeType":     fileUpload.MimeType,
				"uploadType":   fileUpload.UploadType,
				"entityID":     fileUpload.EntityID,
				"createdAt":    fileUpload.CreatedAt,
			})
		}

		c.JSON(http.StatusOK, files)
	}
}

// @Summary Удалить файл
// @Description Удаляет загруженный файл (только владелец файла может его удалить)
// @Tags Files
// @Produce json
// @Param file_id path string true "File ID (UUID)"
// @Security OAuth2Password
// @Success 200 {object} models.SuccessResponse "Success message"
// @Failure 400 {object} models.ErrorResponse "Invalid file ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Access denied"
// @Failure 404 {object} models.ErrorResponse "File not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/api/delete_file/{file_id} [delete]
func DeleteFile(storage *database.Storage, fileService *file.FileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID пользователя из контекста
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Преобразуем строку в UUID
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Получаем ID файла
		fileIDStr := c.Param("file_id")
		fileID, err := uuid.Parse(fileIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
			return
		}

		// Получаем информацию о файле
		ctx := c.Request.Context()
		fileUpload, err := storage.GetFileUploadByID(ctx, fileID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		// Проверяем, что файл принадлежит пользователю
		if fileUpload.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		// Удаляем файл с диска
		if err := fileService.DeleteFile(fileUpload.FilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from disk"})
			return
		}

		// Обновляем связанные таблицы в зависимости от типа файла
		switch fileUpload.UploadType {
		case "profile_image":
			// Очищаем photo_url в таблице пользователей
			if err := storage.UpdateUserProfileImage(ctx, userID, ""); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
				return
			}
		case "resume":
			// Очищаем resume_url в таблице пользователей
			if err := storage.UpdateUserResumeURL(ctx, userID, ""); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user resume"})
				return
			}
		case "achievement_image":
			// Очищаем image_url в таблице достижений
			if fileUpload.EntityID != nil {
				if err := storage.UpdateAchievementImage(ctx, *fileUpload.EntityID, ""); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievement"})
					return
				}
			}
		}

		// Удаляем запись из базы данных
		if err := storage.DeleteFileUpload(ctx, fileID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file record"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
	}
}

// @Summary Загрузить резюме пользователя
// @Description Загружает резюме для текущего пользователя и обновляет поле resume_url
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param resume formData file true "Resume file (PDF, DOC, DOCX, max 10MB)"
// @Security OAuth2Password
// @Success 200 {object} models.FileUploadResponse "Success message with file info"
// @Failure 400 {object} models.ErrorResponse "Invalid file or request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/api/upload_resume [post]
func UploadResume(storage *database.Storage, fileService *file.FileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID пользователя из контекста
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Преобразуем строку в UUID
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Получаем файл из запроса
		fileHeader, err := c.FormFile("resume")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No resume file provided"})
			return
		}

		// Проверяем, что файл является документом
		if err := fileService.ValidateDocumentFile(fileHeader); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Загружаем файл
		fileUpload, err := fileService.UploadFile(fileHeader, userID, "resume", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
			return
		}

		// Сохраняем информацию о файле в базе данных
		ctx := c.Request.Context()
		if err := storage.SaveFileUpload(ctx, fileUpload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file information"})
			return
		}

		// Обновляем URL резюме пользователя
		resumeURL := fileService.GetFileURL(fileUpload.FileName)
		if err := storage.UpdateUserResumeURL(ctx, userID, resumeURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resume URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Resume uploaded successfully",
			"file": gin.H{
				"id":           fileUpload.ID,
				"fileName":     fileUpload.FileName,
				"originalName": fileUpload.OriginalName,
				"fileSize":     fileUpload.FileSize,
				"mimeType":     fileUpload.MimeType,
				"url":          resumeURL,
			},
		})
	}
} 