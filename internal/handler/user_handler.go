package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/AlexRijikov/go-petshop-api/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	user, err := h.repo.GetByID(uint(userID))
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

// GetAllUsers отримує список усіх користувачів (для адміністраторів)

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // встановлюємо таймаут для контексту
	defer cancel()                                                          // забезпечуємо скасування контексту після завершення функції

	users, err := h.repo.GetAll(ctx) // отримуємо всіх користувачів з бази даних
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"}) // якщо сталася помилка при отриманні користувачів, повертаємо 500 Internal Server Error
		return
	}

	c.JSON(http.StatusOK, users) // повертаємо список користувачів у відповіді
}

// DeleteUser видаляє користувача за його ID (для адміністраторів)

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // якщо ID некоректний, повертаємо 400 Bad Request
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // встановлюємо таймаут для контексту
	defer cancel()                                                          // забезпечуємо скасування контексту після завершення функції

	// Викликаємо метод репозиторію для видалення користувача за його ID

	if err := h.repo.Delete(ctx, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"}) // якщо сталася помилка при видаленні, повертаємо 500 Internal Server Error
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"}) // повідомлення про успішне видалення
}

// ChangePassword — зміна паролю користувача

func (h *UserHandler) ChangePassword(c *gin.Context) {
    userID := c.GetInt("user_id") // Отримуємо ID користувача з JWT middleware
    if userID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
        return
    }

    // Структура для отримання старого і нового паролю
    var req struct {
        OldPassword string `json:"old_password"`
        NewPassword string `json:"new_password"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Отримуємо користувача з бази
    user, err := h.repo.GetByID(uint(userID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Перевіряємо, чи старий пароль збігається
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect old password"})
        return
    }

    // Хешуємо новий пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    // Оновлюємо пароль у базі
    if err := h.repo.UpdatePassword(uint(userID), string(hashedPassword)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}


