-- Добавляем поле для изображения достижения
ALTER TABLE achievements ADD COLUMN image_url VARCHAR(500);

-- Создаем таблицу для хранения информации о загруженных файлах
CREATE TABLE file_uploads (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    upload_type VARCHAR(50) NOT NULL, -- 'profile_image', 'achievement_image', etc.
    entity_id UUID, -- ID связанной сущности (achievement_id, user_id, etc.)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индекс для быстрого поиска файлов по пользователю
CREATE INDEX idx_file_uploads_user_id ON file_uploads(user_id);

-- Создаем индекс для быстрого поиска файлов по типу и сущности
CREATE INDEX idx_file_uploads_type_entity ON file_uploads(upload_type, entity_id); 