# ITaM Auth API

API для системы аутентификации и управления пользователями ITaM. Включает функциональность для работы с пользователями, достижениями, запросами, уведомлениями и загрузкой файлов.

## Возможности

- 🔐 Аутентификация и авторизация пользователей
- 👤 Управление профилями пользователей
- 🏆 Система достижений
- 📝 Управление запросами
- 🔔 Система уведомлений
- 📁 Загрузка и управление изображениями
- 📚 Swagger документация

## Быстрый старт

### Требования

- Go 1.19+
- PostgreSQL 12+
- Docker (опционально)

### Установка

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd itam_auth
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте файл `.env` на основе `env.sample`:
```bash
cp env.sample .env
```

4. Настройте переменные окружения в `.env`:
```env
# Database Configuration
DB_USER=itam_user
DB_PASSWORD=itam_db
DB_HOST=localhost
DB_PORT=5432
DB_NAME=itam_auth

# JWT Configuration
JWT_SECRET_KEY=your-secret-key-here

# Migrations
MIGRATIONS_PATH=./migrations

# File Upload Configuration
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760
ALLOWED_TYPES=.jpg,.jpeg,.png,.gif,.webp
```

5. Примените миграции:
```bash
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
```

6. Запустите сервер:
```bash
go run cmd/app/main.go
```

Сервер будет доступен по адресу: `http://localhost:8080`

## API Документация

Swagger документация доступна по адресу: `http://localhost:8080/auth/swagger`

### Основные эндпоинты

#### Аутентификация
- `POST /auth/api/register` - Регистрация пользователя
- `POST /auth/api/login` - Авторизация пользователя

#### Пользователи
- `GET /auth/api/me` - Получить текущего пользователя
- `PATCH /auth/api/update_user_info` - Обновить информацию пользователя
- `GET /auth/api/get_user/{user_id}` - Получить пользователя по ID

#### Достижения
- `POST /auth/api/create_achievement` - Создать достижение
- `GET /auth/api/get_all_achievements` - Получить все достижения
- `GET /auth/api/get_user_achievements` - Получить достижения пользователя

#### Файлы
- `POST /auth/api/upload_profile_image` - Загрузить изображение профиля
- `POST /auth/api/upload_achievement_image` - Загрузить изображение достижения
- `GET /uploads/{filename}` - Получить файл
- `DELETE /auth/api/delete_file/{file_id}` - Удалить файл

## Загрузка файлов

Система поддерживает загрузку изображений для профилей пользователей и достижений.

### Поддерживаемые форматы
- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- WebP (.webp)

### Ограничения
- Максимальный размер файла: 10MB (настраивается)
- Только изображения
- Файлы сохраняются с уникальными именами

### Примеры использования

#### Загрузка изображения профиля
```bash
curl -X POST http://localhost:8080/auth/api/upload_profile_image \
  -H "Authorization: Bearer <your-jwt-token>" \
  -F "image=@profile.jpg"
```

#### Загрузка изображения достижения
```bash
curl -X POST http://localhost:8080/auth/api/upload_achievement_image \
  -H "Authorization: Bearer <your-jwt-token>" \
  -F "image=@achievement.jpg" \
  -F "achievement_id=uuid-of-achievement"
```

## Структура проекта

```
itam_auth/
├── cmd/
│   ├── app/
│   │   └── main.go
│   └── migrator/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   ├── services/
│   └── utils/
├── migrations/
├── docs/
└── uploads/
```

## Разработка

### Запуск в режиме разработки
```bash
go run cmd/app/main.go
```

### Запуск миграций
```bash
# Применение миграций
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up

# Откат миграций
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down
```

### Обновление Swagger документации
```bash
swag init -g cmd/app/main.go -o docs
```

## Docker

### Сборка образа
```bash
docker build -t itam-auth .
```

### Запуск контейнера
```bash
docker run -p 8080:8080 itam-auth
```

### Docker Compose
```bash
docker-compose up -d
```

## Безопасность

- JWT токены для аутентификации
- Валидация входных данных
- Проверка типов файлов
- Ограничение размера файлов
- Проверка прав доступа

## Лицензия

MIT License

## Поддержка

- Email: support@itam.com
- Документация: `http://localhost:8080/auth/swagger`
