package repositories

import (
	"context"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"gorm.io/gorm"
)

// List повертає і загальну кількість (потрібно для пагінації на фронті).

// Всі методи використовують WithContext(ctx) — корисно для таймаутів/тестів.

type ProductRepository interface {
	Create(ctx context.Context, p *models.Product) error                          // p.ID заповнюється автоматично
	GetByID(ctx context.Context, id uint) (*models.Product, error)                // повертає nil, nil якщо не знайдено
	List(ctx context.Context, limit, offset int) ([]models.Product, int64, error) // returns items, totalCount
	Update(ctx context.Context, p *models.Product) error
	Delete(ctx context.Context, id uint) error
}

// productRepo реалізує ProductRepository

type productRepo struct {
	db *gorm.DB
}

// NewProductRepository створює новий ProductRepository

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

// Create додає новий продукт в базу даних

func (r *productRepo) Create(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// GetByID шукає продукт за ID

func (r *productRepo) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var p models.Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// List повертає продукти з пагінацією

func (r *productRepo) List(ctx context.Context, limit, offset int) ([]models.Product, int64, error) {
	var items []models.Product
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Product{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Update змінює дані продукту

func (r *productRepo) Update(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

// Delete видаляє продукт за ID

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
