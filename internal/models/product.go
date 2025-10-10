package models

import (
	"time"

	"gorm.io/gorm"
)

// PriceCents — зберігаємо в цілих (копійках) щоб уникнути FP-помилок.

// DeletedAt для soft-delete (індекс).

// JSON-теги для відповіді API.

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`                      // Primary key (Первинний ключ)
	CreatedAt   time.Time      `json:"created_at"`                                // Час створення запису
	UpdatedAt   time.Time      `json:"updated_at"`                                // Час останнього оновлення запису
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                            // Soft delete (м'яке видалення)
	Name        string         `gorm:"size:255;not null" json:"name"`             // Назва продукту
	Description string         `gorm:"type:text" json:"description,omitempty"`    // Опис продукту
	PriceCents  int64          `gorm:"not null" json:"price_cents"`               // ціна в центі (копійки) (щоб уникнути float)
	Stock       int            `gorm:"not null;default:0" json:"stock"`           // Кількість на складі
	SKU         string         `gorm:"size:100;uniqueIndex" json:"sku,omitempty"` // Унікальний артикул (Stock Keeping Unit)
}
