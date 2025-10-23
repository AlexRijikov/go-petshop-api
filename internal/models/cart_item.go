package models

import (
	"time"

	"gorm.io/gorm"
)

// CartItem представляє товар, доданий до кошика користувача.
// UserID та ProductID є зовнішніми ключами, що посилаються на відповідні таблиці.
// JSON-теги використовуються для відповіді API.
// DeletedAt для soft-delete (індекс).	
// UserID та ProductID мають індекси для швидкого пошуку товарів у кошику конкретного користувача.
// Quantity має значення за замовчуванням 1, оскільки при додаванні товару до кошика зазвичай додається один екземпляр.
// DeletedAt має індекс для ефективного виконання запитів, що виключають видалені записи.


type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`               // Primary key (Первинний ключ)
	CreatedAt time.Time      `json:"created_at"`                         // Час створення запису
	UpdatedAt time.Time      `json:"updated_at"`                         // Час останнього оновлення запису
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                     // Soft delete (м'яке видалення)
	UserID    uint           `gorm:"not null;index" json:"user_id"`      // Ідентифікатор користувача чий саме кошик (зовнішній ключ) 
	ProductID uint           `gorm:"not null;index" json:"product_id"`   // Ідентифікатор продукту (зовнішній ключ)
	Quantity  int            `gorm:"not null;default:1" json:"quantity"` // Кількість товару в кошику
}
