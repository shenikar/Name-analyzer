package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/shenikar/Name-analyzer/config"
	"github.com/shenikar/Name-analyzer/internal/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}
	db, err := db.ConnDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to create db: %v", err)
	}
	defer db.Conn.Close()

	// Проверяем подключение к БД
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// // Пробуем выполнить простой запрос
	// var result int
	// err = db.Conn.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	// if err != nil {
	// 	log.Fatalf("failed to ping database: %v", err)
	// }

	log.Println("Successfully connected to database!")
}
