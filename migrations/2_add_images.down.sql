-- Удаляем таблицу загруженных файлов
DROP TABLE IF EXISTS file_uploads;

-- Удаляем поле изображения достижения
ALTER TABLE achievements DROP COLUMN IF EXISTS image_url; 