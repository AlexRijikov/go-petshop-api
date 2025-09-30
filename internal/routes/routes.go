package routes

import (
	"github.com/AlexRijikov/go-petshop-api/internal/handler"
	"github.com/AlexRijikov/go-petshop-api/internal/repository"
	"github.com/AlexRijikov/go-petshop-api/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes відповідає за підключення всіх маршрутів API

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	// Health-check маршрут

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Репозиторії (робота напряму з БД)

	productRepo := repositories.NewProductRepository(db)

	// Сервіси (бізнес-логіка)

	productService := services.NewProductService(productRepo)

	// Хендлери (обробка HTTP-запитів)

	productHandler := handlers.NewProductHandler(productService)

	// Група маршрутів /api/products

	api := r.Group("/api")
	productHandler.RegisterRoutes(api)
}
