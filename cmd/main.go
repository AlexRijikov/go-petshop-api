package main

import (
	"log"

	"github.com/AlexRijikov/go-petshop-api/internal/database"
	"github.com/AlexRijikov/go-petshop-api/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	// підключення до БД (з автоматичною міграцією)

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// створення Gin-сервера

	r := gin.Default()

	// реєстрація маршрутів

	routes.RegisterRoutes(r, db)

	log.Println(" Сервер запущено на local host :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
