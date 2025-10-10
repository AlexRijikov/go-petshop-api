package services

import (
	"context"
	"errors"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"github.com/AlexRijikov/go-petshop-api/internal/repository"
)

// Помилки сервісу продуктів

var (
	ErrInvalidPrice = errors.New("price must be > 0") // ціна має бути більшою за 0
	ErrNotFound     = errors.New("product not found") // продукт не знайдено
)

// ProductService визначає бізнес-логіку для продуктів

type ProductService interface {
	CreateProduct(ctx context.Context, p *models.Product) (*models.Product, error)        // p.ID заповнюється автоматично
	GetProduct(ctx context.Context, id uint) (*models.Product, error)                     // повертає ErrNotFound якщо не знайдено або іншу помилку
	ListProducts(ctx context.Context, limit, offset int) ([]models.Product, int64, error) // returns items, totalCount
	UpdateProduct(ctx context.Context, p *models.Product) (*models.Product, error)        // повертає ErrNotFound якщо не знайдено або ErrInvalidPrice якщо ціна некоректна
	DeleteProduct(ctx context.Context, id uint) error                                     // повертає ErrNotFound якщо не знайдено
}

// productService реалізує ProductService

type productService struct {
	repo repositories.ProductRepository
}

// NewProductService створює новий ProductService

func NewProductService(r repositories.ProductRepository) ProductService {
	return &productService{repo: r}
}

// CreateProduct створює новий продукт, перевіряє що ціна > 0 (в копійках)

func (s *productService) CreateProduct(ctx context.Context, p *models.Product) (*models.Product, error) {
	if p.PriceCents <= 0 {
		return nil, ErrInvalidPrice
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// GetProduct повертає продукт за ID або ErrNotFound якщо не знайдено (репозиторій повертає помилку)

func (s *productService) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return p, nil
}

// ListProducts повертає продукти з пагінацією (limit, offset)

func (s *productService) ListProducts(ctx context.Context, limit, offset int) ([]models.Product, int64, error) {
	return s.repo.List(ctx, limit, offset)
}

// UpdateProduct оновлює продукт, перевіряє що ціна > 0 (в копійках)

func (s *productService) UpdateProduct(ctx context.Context, p *models.Product) (*models.Product, error) {
	if p.PriceCents <= 0 {
		return nil, ErrInvalidPrice
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// DeleteProduct видаляє продукт за ID (повертає ErrNotFound якщо не знайдено)

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
