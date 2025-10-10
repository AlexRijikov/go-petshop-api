package routes

import (
	"github.com/AlexRijikov/go-petshop-api/internal/handler"
	"github.com/AlexRijikov/go-petshop-api/internal/middleware"
	"github.com/AlexRijikov/go-petshop-api/internal/repository"
	"github.com/AlexRijikov/go-petshop-api/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes реєструє всі маршрути (ендпоінти) для продуктів та аутентифікації

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// Групуємо всі маршрути під префіксом /api
	api := r.Group("/api")

	// Реєструємо маршрути для продуктів

	productRepo := repositories.NewProductRepository(db)     // створюємо репозиторій продуктів
	productSvc := services.NewProductService(productRepo)    // створюємо сервіс продуктів
	productHandler := handlers.NewProductHandler(productSvc) // створюємо хендлер продуктів
	productHandler.RegisterRoutes(api)

	// Реєструємо маршрути для аутентифікації

	userRepo := repositories.NewUserRepository(db)  // створюємо репозиторій користувачів
	authSvc := services.NewAuthService(userRepo)    // створюємо сервіс аутентифікації(реєстрація/вхід)
	authHandler := handlers.NewAuthHandler(authSvc) // створюємо хендлер аутентифікації
	authHandler.RegisterRoutes(api)                 // реєструємо маршрути аутентифікації

	// Тестовий ендпоінт для перевірки, що сервер працює (можна видалити пізніше)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Реєструємо маршрути для користувачів (профіль, оновлення профілю) з захистом за допомогою JWT токенів через AuthMiddleware

	authMiddleware := middleware.AuthMiddleware()    // створюємо middleware для аутентифікації користувачів
	userHandler := handlers.NewUserHandler(userRepo) // створюємо хендлер користувачів

	protected := api.Group("/users") // групуємо маршрути користувачів під префіксом /api/users
	protected.Use(authMiddleware)    // застосовуємо middleware для захисту маршрутів
	{
		protected.GET("/me", userHandler.GetProfile)    // отримання профілю користувача
		protected.PUT("/me", userHandler.UpdateProfile) // оновлення профілю користувача (наприклад, email або username)

	}
}
