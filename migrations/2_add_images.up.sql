-- Добавляем поле для изображения достижения
ALTER TABLE achievements ADD COLUMN image_url TEXT;

-- Создаем таблицу для хранения информации о загруженных файлах
CREATE TABLE file_uploads (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    file_name TEXT NOT NULL,
    original_name TEXT,
    file_path TEXT,
    file_size BIGINT,
    mime_type TEXT,
    upload_type TEXT, -- 'profile_image', 'achievement_image', etc.
    entity_id UUID, -- ID связанной сущности (achievement_id, user_id, etc.)
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создаем индекс для быстрого поиска файлов по пользователю
CREATE INDEX idx_file_uploads_user_id ON file_uploads(user_id);

-- Создаем индекс для быстрого поиска файлов по типу и сущности
CREATE INDEX idx_file_uploads_type_entity ON file_uploads(upload_type, entity_id); 