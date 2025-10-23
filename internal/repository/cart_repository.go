package repositories

import (
	"context"
	"gorm.io/gorm"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
)

// CartRepository визначає інтерфейс для роботи з кошиком користувача

type CartRepository interface {
	ListByUserID(ctx context.Context, userID uint) ([]models.CartItem, error)
	Add(ctx context.Context, item *models.CartItem) error
	UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error
	Remove(ctx context.Context, userID, productID uint) error
}

// cartRepository реалізує інтерфейс CartRepository

type cartRepository struct {
	db *gorm.DB // з'єднання з базою даних
}

// NewCartRepository створює новий CartRepository з наданим з'єднанням з базою даних

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db} // повертаємо новий екземпляр cartRepository
}

// ListByUserID отримує всі товари в кошику користувача за його ID

func (r *cartRepository) ListByUserID(ctx context.Context, userID uint) ([]models.CartItem, error) {
	var items []models.CartItem // створюємо зріз для зберігання товарів кошика
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err // повертаємо помилку, якщо виникла проблема з отриманням даних
	}
	return items, nil // повертаємо отримані товари кошика
}

// Add - 

func (r *cartRepository) Add(ctx context.Context, item *models.CartItem) error {
	// якщо товар вже є — просто оновимо кількість
	var existing models.CartItem
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).
		First(&existing).Error

	if err == nil {
		existing.Quantity += item.Quantity
		return r.db.WithContext(ctx).Save(&existing).Error
	}
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(item).Error
	}
	return err
}

// UpdateQuantity - 

func (r *cartRepository) UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Update("quantity", quantity).Error
}

// Remove -

func (r *cartRepository) Remove(ctx context.Context, userID, productID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.CartItem{}).Error
}
