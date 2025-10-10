package handlers

import (
	"net/http"

	"github.com/AlexRijikov/go-petshop-api/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler обробляє HTTP-запити, пов'язані з аутентифікацією користувачів (реєстрація та вхід)

type AuthHandler struct {
	svc services.AuthService
}

// NewAuthHandler створює новий AuthHandler з наданим сервісом AuthService

func NewAuthHandler(s services.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

// RegisterRoutes реєструє маршрути аутентифікації у вказаній групі маршрутизатора (rg *gin.RouterGroup)

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
}

// authRequest використовується для прив'язки та валідації вхідних даних при реєстрації та вході користувача

type authRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register обробляє реєстрацію нового користувача

func (h *AuthHandler) Register(c *gin.Context) {

	// Прив'язуємо та валідовуємо вхідні дані
	var req authRequest

	// Якщо помилка прив'язки/валідації, повертаємо 400 Bad Request (помилка клієнта)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Викликаємо сервіс для реєстрації користувача

	if err := h.svc.Register(c.Request.Context(), req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Повертаємо успішну  відповідь при успішній реєстрації (статус 201 Created)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

// Login обробляє вхід користувача і повертає JWT токен при успішній аутентифікації

func (h *AuthHandler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Повертаємо токен у відповіді

	c.JSON(http.StatusOK, gin.H{"token": token})
}
