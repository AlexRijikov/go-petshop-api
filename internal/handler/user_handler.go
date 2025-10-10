package handlers

import (
	"net/http"
	
	"github.com/AlexRijikov/go-petshop-api/internal/repository" 
	"github.com/gin-gonic/gin"
)

// UserHandler відповідає за обробку HTTP-запитів, пов'язаних із користувачами (отримання профілю, оновлення профілю)

type UserHandler struct {
	repo repositories.UserRepository
}

// NewUserHandler створює новий екземпляр UserHandler з наданим репозиторієм користувачів (repository.UserRepository)

func NewUserHandler(r repositories.UserRepository) *UserHandler {
	return &UserHandler{repo: r}
}

// GetProfile — отримання профілю користувача

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id") // отримуємо user_id з контексту, встановленого AuthMiddleware

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"}) // якщо user_id відсутній, повертаємо 401 Unauthorized
		return
	}

	// Отримуємо дані користувача з бази даних за допомогою репозиторію

	user, err := h.repo.GetByID (uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // якщо користувача не знайдено, повертаємо 404 Not Found
		return
	}

	// Повертаємо дані користувача у відповіді (без пароля)

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

// UpdateProfile — оновлення даних користувача (наприклад, email або username)

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("user_id") // отримуємо user_id з контексту, встановленого AuthMiddleware
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"}) // якщо user_id відсутній, повертаємо 401 Unauthorized
		return
	}

	// Прив'язуємо вхідні дані (username, email) з JSON тіла запиту

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // якщо помилка прив'язки, повертаємо 400 Bad Request
		return
	}

	// Оновлюємо дані користувача в базі даних за допомогою репозиторію

	user, err := h.repo.UpdateProfile(uint(userID), req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"}) // якщо сталася помилка при оновленні, повертаємо 500 Internal Server Error
		return
	}

	// Повертаємо оновлені дані користувача у відповіді (без пароля)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Profile updated successfully", // повідомлення про успішне оновлення
		"user_id":  user.ID,                        // повертаємо ID користувача
		"username": user.Username,                  // повертаємо оновлене ім'я користувача
		"email":    user.Email,                     // повертаємо оновлений email
		"role":     user.Role,                      // повертаємо роль користувача (наприклад, "user" або "admin")
	})
}
