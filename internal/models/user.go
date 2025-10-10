package models

import (
	"gorm.io/gorm"
	"time"
)

// User представляє користувача системи.
// Пароль зберігається в хешованому вигляді.
// Роль визначає рівень доступу користувача (наприклад, "user", "admin").
// JSON-теги використовуються для відповіді API.
// DeletedAt для soft-delete (індекс).

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`                          // Primary key (Первинний ключ)
	CreatedAt time.Time      `json:"created_at"`                                    // Час створення запису
	UpdatedAt time.Time      `json:"updated_at"`                                    // Час останнього оновлення запису
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                // Soft delete (м'яке видалення)
	Username  string         `gorm:"size:255;not null;uniqueIndex" json:"username"` // Ім'я користувача
	Email     string         `gorm:"size:255;not null;uniqueIndex" json:"email"`    // Електронна пошта користувача
	Password  string         `gorm:"size:255;not null" json:"-"`                    // Хешований пароль (не включається в JSON-відповідь)
	Role      string         `gorm:"size:50;not null;default:'user'" json:"role"`   // Ролі можуть бути : user, admin
}
