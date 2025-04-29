package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/shenikar/Name-analyzer/config"
	"github.com/shenikar/Name-analyzer/internal/api"
	"github.com/shenikar/Name-analyzer/internal/db"
)

// @title Name Analyzer API
// @version 1.0
// @description Сервис для анализа имен и обогащения данных о людях (возраст, пол, национальность)

// @host localhost:8080
// @BasePath /api/v1

// @contact.name API Support
// @contact.email your-email@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	// Создаем конфигурацию приложения из переменных окружения
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("не удалось создать конфигурацию: %v", err)
	}

	// Инициализируем логгер с выводом номеров строк в логах
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// Применяем миграции к базе данных
	if err := db.RunMigrations(cfg.DBDSN); err != nil {
		log.Fatalf("не удалось применить миграции: %v", err)
	}

	// Устанавливаем соединение с базой данных
	db, err := db.ConnDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}

	// Создаем новый роутер
	mux := http.NewServeMux()
	// Регистрируем все API маршруты
	api.RegisterRoutes(mux, db, logger)

	// Добавляем промежуточное ПО (middleware) для логирования и обработки паник
	handler := api.LoggingMiddleware(logger)(api.RecoverMiddleware(mux))

	// Запускаем HTTP сервер на указанном порту
	logger.Printf("Server is running on port: %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
