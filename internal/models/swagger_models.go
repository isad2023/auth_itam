package models

// FileUploadResponse представляет ответ при загрузке файла
type FileUploadResponse struct {
	Message string `json:"message" example:"File uploaded successfully"`
	File    FileInfo `json:"file"`
}

// FileInfo представляет информацию о загруженном файле
type FileInfo struct {
	ID           string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FileName     string `json:"fileName" example:"550e8400-e29b-41d4-a716-446655440000_1234567890.jpg"`
	OriginalName string `json:"originalName" example:"profile.jpg"`
	FileSize     int64  `json:"fileSize" example:"12345"`
	MimeType     string `json:"mimeType" example:"image/jpeg"`
	URL          string `json:"url" example:"/uploads/550e8400-e29b-41d4-a716-446655440000_1234567890.jpg"`
}

// ErrorResponse представляет стандартный ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid file format"`
	Details string `json:"details,omitempty" example:"Only JPEG, PNG, GIF, WebP formats are allowed"`
}

// SuccessResponse представляет стандартный ответ об успехе
type SuccessResponse struct {
	Message string `json:"message" example:"File deleted successfully"`
}

// LoginResponse представляет ответ на успешную авторизацию
type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"2592000"`
}

// RegisterResponse представляет ответ на успешную регистрацию
type RegisterResponse struct {
	Message string `json:"message" example:"User registered successfully"`
	User    User   `json:"user"`
} 