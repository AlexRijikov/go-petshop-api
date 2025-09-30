package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {

	// Завантажуємо .env
	
	err := godotenv.Load()
	if err != nil {
		log.Println(" Не вдалося завантажити .env, використовую системні змінні")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Формуємо DSN

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	// Підключення

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Не вдалося підключитися до БД: %w", err)
	}

	// Автоматична міграція моделей

	if err := db.AutoMigrate(&models.Product{}); err != nil {
		return nil, fmt.Errorf("Помилка AutoMigrate: %w", err)
	}

	DB = db
	fmt.Println("Успішне підключення та міграція PostgreSQL")
	return db, nil
}
