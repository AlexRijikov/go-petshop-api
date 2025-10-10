package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Секретний ключ для підпису JWT токенів - (у реальному застосунку зберігайте його в безпечному місці )

var jwtKey = []byte("supersecretkey") // потім винесемо у .env

// AuthMiddleware перевіряє JWT токен в заголовку Authorization
// Якщо токен дійсний, додає user_id в контекст запиту

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Очікуємо формат: Bearer <token>

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Парсимо токен і перевіряємо його дійсність (підпис, термін дії тощо)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		// Якщо токен недійсний або сталася помилка, повертаємо 401 Unauthorized

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Отримуємо claims і додаємо user_id в контекст запиту для подальшого використання в обробниках запитів (handlers)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		}

		c.Next()
	}
}
