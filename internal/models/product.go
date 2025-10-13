package models

import (
	"time"

	"gorm.io/gorm"
)

// PriceCents — зберігаємо в цілих (копійках) щоб уникнути FP-помилок.

// DeletedAt для soft-delete (індекс).

// JSON-теги для відповіді API.

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`                      // Primary key (Первинний ключ - унікальний ідентифікатор продукту в базі даних для швидкого пошуку та зв'язку з іншими таблицями)
	CreatedAt   time.Time      `json:"created_at"`                                // Час створення запису (створення продукту  в системі для відстеження коли продукт був доданий до системи)
	UpdatedAt   time.Time      `json:"updated_at"`                                // Час останнього оновлення запису 
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                            // Soft delete (м'яке видалення з індексом для швидкого пошуку не видалених записів)
	Name        string         `gorm:"size:255;not null" json:"name"`             // Назва продукту (* обов'язково - для ідентифікації продукту в системі та для відображення користувачам)
	Description string         `gorm:"type:text" json:"description,omitempty"`    // Опис продукту (опціонально - для детального опису продукту )
	PriceCents  int64          `gorm:"not null" json:"price_cents"`               // ціна в центі (копійки) (щоб уникнути float)
	Stock       int            `gorm:"not null;default:0" json:"stock"`           // Кількість на складі (Stock - для відстеження кількості продуктів на складі)
	SKU         string         `gorm:"size:100;uniqueIndex" json:"sku,omitempty"` // Унікальний артикул (Stock Keeping Unit - для відстеження запасів продуктів  в системі управління запасами або ERP  системі  наприклад SAP, Oracle і т.д.)
	ImageURL    string         `gorm:"size:255" json:"image_url,omitempty"`       // URL зображення продукту (опціонально - для відображення зображення продукту )
	Category    string         `gorm:"size:100" json:"category,omitempty"`        // Категорія продукту (опціонально- для фільтрації та сортування  продуктів за категоріями наприклад корм, сушені смаколики і т.д. )
	Metadata    string         `gorm:"type:json" json:"metadata,omitempty"`       // Додаткові метадані у форматі JSON (опціонально - для розширення інформації про продукт наприклад колір, розмір і т.д.)
	




}
