package models

import (
	"time"

	"gorm.io/gorm"
)

// PriceCents — зберігаємо в цілих (копійках) щоб уникнути FP-помилок.

// DeletedAt для soft-delete (індекс).

// JSON-теги для відповіді API.

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	PriceCents  int64          `gorm:"not null" json:"price_cents"` // ціна в центі (копійки) (щоб уникнути float)
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	SKU         string         `gorm:"size:100;uniqueIndex" json:"sku,omitempty"`
}
