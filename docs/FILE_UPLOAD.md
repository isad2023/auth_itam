# Система загрузки изображений

## Обзор

Система позволяет загружать изображения для профилей пользователей (используя существующее поле `photo_url`) и достижений. Все файлы сохраняются локально в указанной директории и информация о них записывается в базу данных.

## Конфигурация

Добавьте следующие переменные окружения в ваш `.env` файл:

```env
# Путь для сохранения загруженных файлов
UPLOAD_PATH=./uploads

# Максимальный размер файла в байтах (по умолчанию 10MB)
MAX_FILE_SIZE=10485760

# Разрешенные типы файлов (через запятую)
ALLOWED_TYPES=.jpg,.jpeg,.png,.gif,.webp
```

## API Endpoints

### Загрузка изображения профиля

**POST** `/auth/api/upload_profile_image`

Загружает изображение профиля для текущего пользователя (обновляет поле `photo_url`).

**Параметры:**
- `image` (file) - файл изображения

**Заголовки:**
- `Authorization: Bearer <token>` - JWT токен

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/auth/api/upload_profile_image \
  -H "Authorization: Bearer <your-jwt-token>" \
  -F "image=@profile.jpg"
```

**Ответ:**
```json
{
  "message": "Profile image uploaded successfully",
  "file": {
    "id": "uuid",
    "fileName": "generated-filename.jpg",
    "originalName": "profile.jpg",
    "fileSize": 12345,
    "mimeType": "image/jpeg",
    "url": "/uploads/generated-filename.jpg"
  }
}
```

### Загрузка изображения достижения

**POST** `/auth/api/upload_achievement_image`

Загружает изображение для достижения (обновляет поле `image_url`).

**Параметры:**
- `image` (file) - файл изображения
- `achievement_id` (string) - ID достижения

**Заголовки:**
- `Authorization: Bearer <token>` - JWT токен

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/auth/api/upload_achievement_image \
  -H "Authorization: Bearer <your-jwt-token>" \
  -F "image=@achievement.jpg" \
  -F "achievement_id=uuid-of-achievement"
```

### Получение файла

**GET** `/uploads/{filename}`

Возвращает файл по имени.

**Пример:**
```
GET /uploads/generated-filename.jpg
```

### Удаление файла

**DELETE** `/auth/api/delete_file/{file_id}`

Удаляет загруженный файл.

**Заголовки:**
- `Authorization: Bearer <token>` - JWT токен

**Пример запроса:**
```bash
curl -X DELETE http://localhost:8080/auth/api/delete_file/uuid-of-file \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Поддерживаемые форматы

- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- WebP (.webp)

## Ограничения

- Максимальный размер файла: 10MB (настраивается)
- Только изображения
- Файлы сохраняются с уникальными именами для предотвращения конфликтов

## Структура базы данных

### Таблица `file_uploads`

```sql
CREATE TABLE file_uploads (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    upload_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Обновленные таблицы

**users:**
- `photo_url VARCHAR(255)` - URL изображения профиля (существующее поле)

**achievements:**
- `image_url VARCHAR(500)` - URL изображения достижения (новое поле)

## Безопасность

- Проверка MIME типов
- Валидация расширений файлов
- Ограничение размера файлов
- Проверка прав доступа (только владелец может удалять файлы)
- Безопасные имена файлов (UUID + timestamp)

## Миграции

Для применения изменений выполните миграцию:

```bash
# Применение миграции
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up

# Откат миграции
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down
```

## Интеграция с существующим кодом

Система интегрируется с существующими моделями:

- **User.PhotoURL** - используется для хранения URL изображения профиля
- **Achievement.ImageURL** - новое поле для хранения URL изображения достижения

При загрузке изображения профиля обновляется поле `photo_url` в таблице `users`, что обеспечивает обратную совместимость с существующим кодом. 