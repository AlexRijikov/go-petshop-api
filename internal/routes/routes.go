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

	// AUTH - маршрути для реєстрації, входу, виходу  (реєстрація, логін) — публічні маршрути (без авторизації)

	userRepo := repositories.NewUserRepository(db)  // створюємо репозиторій користувачів
	authSvc := services.NewAuthService(userRepo)    // створюємо сервіс аутентифікації з репозиторієм користувачів
	authHandler := handlers.NewAuthHandler(authSvc) // створюємо хендлер аутентифікації з сервісом аутентифікації
	authHandler.RegisterRoutes(api)

	// PRODUCTS - маршрути для роботи з товарами ( отримання списку продуктів, створення продукту тощо) — публічні маршрути (без авторизації)

	productRepo := repositories.NewProductRepository(db)     // створюємо репозиторій продуктів
	productSvc := services.NewProductService(productRepo)    // створюємо сервіс продуктів з репозиторієм продуктів
	productHandler := handlers.NewProductHandler(productSvc) // створюємо хендлер продуктів з сервісом продуктів
	productHandler.RegisterRoutes(api)                       // реєструємо маршрути продуктів

	// USERS - отримання профілю, оновлення профілю користувача тощо — захищені маршрути AuthMiddleware (перевірка JWT)

	authMiddleware := middleware.AuthMiddleware()    // створюємо middleware для авторизації (перевірка JWT)
	userHandler := handlers.NewUserHandler(userRepo) // створюємо хендлер користувачів з репозиторієм користувачів

	users := api.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("/me", userHandler.GetProfile)
		users.PUT("/me", userHandler.UpdateProfile)
		users.PUT("/me/password", userHandler.ChangePassword)

	}

	//  Ping endpoint для перевірки стану сервера (можна видалити в продакшені)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}
