package repositories

import (
	"context"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"gorm.io/gorm"
)

// List повертає і загальну кількість (потрібно для пагінації на фронті).

// Всі методи використовують WithContext(ctx) — корисно для таймаутів/тестів.

type ProductRepository interface {
	Create(ctx context.Context, p *models.Product) error
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	List(ctx context.Context, limit, offset int) ([]models.Product, int64, error) // returns items, totalCount
	Update(ctx context.Context, p *models.Product) error
	Delete(ctx context.Context, id uint) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *productRepo) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var p models.Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

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

func (r *productRepo) Update(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
