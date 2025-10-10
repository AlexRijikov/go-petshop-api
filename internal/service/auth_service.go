package services

import (
	"context"
	"errors"
	"time"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"github.com/AlexRijikov/go-petshop-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("supersecretkey") // потім винесемо у .env

// AuthService відповідає за реєстрацію та логін користувачів

type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
}

// authService реалізує AuthService

type authService struct {
	repo repositories.UserRepository
}

// NewAuthService створює новий AuthService

func NewAuthService(r repositories.UserRepository) AuthService {
	return &authService{repo: r}
}

// Register створює нового користувача з хешованим паролем

func (s *authService) Register(ctx context.Context, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Створюємо користувача
	user := &models.User{
		Email:    email,
		Password: string(hashed),
	}
	// Зберігаємо користувача в базу даних
	return s.repo.Create(ctx, user)
}

// Login перевіряє email і пароль, повертає JWT токен якщо успішно увійшли в систему

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	// Перевіряємо пароль

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Створюємо JWT токен з user ID і терміном дії 24 години

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	// Підписуємо токен і повертаємо його

	return token.SignedString(jwtKey)
}
