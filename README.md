# Name Analyzer

Сервис для анализа и обогащения данных о людях. Получает базовую информацию о человеке (имя, фамилия, отчество) и обогащает её данными о предполагаемом возрасте, поле и национальности, используя внешние API.

## Функциональность

- Создание, чтение, обновление и удаление записей о людях
- Автоматическое обогащение данных через внешние API:
  - Возраст (agify.io)
  - Пол (genderize.io)
  - Национальность (nationalize.io)
- Хранение данных в PostgreSQL
- REST API с JSON форматом
- Swagger документация
- Логирование операций
- Конфигурация через переменные окружения

## Технологии

- Go 1.22+
- PostgreSQL
- Docker & Docker Compose
- Внешние API: agify.io, genderize.io, nationalize.io

## Установка и запуск

### Предварительные требования

- Go 1.22 или выше
- Docker и Docker Compose
- Make (опционально)

### Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/shenikar/Name-analyzer.git
cd Name-analyzer
```

2. Создайте файл .env на основе примера:
```bash
cp .env_example .env
```

3. Установите зависимости:
```bash
go mod download
```

### Запуск

1. Запустите PostgreSQL через Docker:
```bash
docker compose up -d
```

2. Запустите приложение:
```bash
go run cmd/main.go
```

Сервис будет доступен по адресу: http://localhost:8080

## Примеры использования API

### Создание записи
```bash
curl -X POST http://localhost:8080/api/v1/persons \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван",
    "surname": "Иванов",
    "patronymic": "Иванович"
  }'
```

### Получение списка
```bash
# Все записи (с лимитом 10)
curl http://localhost:8080/api/v1/persons

# С фильтрацией по имени
curl "http://localhost:8080/api/v1/persons?name=Иван"

# С пагинацией
curl "http://localhost:8080/api/v1/persons?limit=5&offset=0"
```

### Получение записи по ID
```bash
curl http://localhost:8080/api/v1/persons/39755c70-2ddb-4a62-90ea-1eeaf07a545a
```

### Обновление записи
```bash
curl -X PUT http://localhost:8080/api/v1/persons/39755c70-2ddb-4a62-90ea-1eeaf07a545a \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван",
    "surname": "Петров"
  }'
```

### Удаление записи
```bash
curl -X DELETE http://localhost:8080/api/v1/persons/39755c70-2ddb-4a62-90ea-1eeaf07a545a
```

### Примеры ответов

Успешное создание записи:
```json
{
  "id": "39755c70-2ddb-4a62-90ea-1eeaf07a545a",
  "name": "Иван",
  "surname": "Иванов",
  "patronymic": "Иванович",
  "age": 44,
  "gender": "male",
  "nationality": "RU",
  "created_at": "2024-03-20T15:04:05Z",
  "updated_at": "2024-03-20T15:04:05Z"
}
```

Ошибка:
```json
{
  "error": "некорректный запрос"
}
```
```

## Конфигурация

Настройки приложения задаются через переменные окружения:

- `DB_DSN` - строка подключения к PostgreSQL
- `PORT` - порт для HTTP сервера (по умолчанию 8080)
- `LOG_LEVEL` - уровень логирования
- `RATE_LIMIT` - ограничение запросов в секунду

## Структура проекта
├── cmd/
│ └── main.go # Точка входа приложения
├── internal/
│ ├── api/ # Обработчики HTTP запросов
│ ├── db/ # Работа с базой данных
│ ├── enrich/ # Обогащение данных через внешние API
│ └── model/ # Модели данных
├── migrations/ # SQL миграции
├── config/ # Конфигурация приложения
├── docs/ # Swagger документация
├── docker-compose.yml # Конфигурация Docker
└── .env # Переменные окружения

## Разработка

### Генерация Swagger документации

```bash
swag init -g cmd/main.go
```

Swagger UI будет доступен по адресу: http://localhost:8080/swagger/index.html

### Миграции

Создание новой миграции:
```bash
migrate create -ext sql -dir migrations -seq create_persons_table
```

Применение миграций:
```bash
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
```

## Лицензия

[MIT License](LICENSE)