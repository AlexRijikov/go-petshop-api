package services

import (
	"context"
	"errors"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"github.com/AlexRijikov/go-petshop-api/internal/repository"
)

//	Сервіс робить базову валідацію (ціна >0), інші правила можна доповнювати.
//
// Повертаємо спеціальні помилки (ErrInvalidPrice, ErrNotFound) — зручніше тестувати й мапити на HTTP-коди в handler.
var (
	ErrInvalidPrice = errors.New("price must be > 0")
	ErrNotFound     = errors.New("product not found")
)

type ProductService interface {
	CreateProduct(ctx context.Context, p *models.Product) (*models.Product, error)
	GetProduct(ctx context.Context, id uint) (*models.Product, error)
	ListProducts(ctx context.Context, limit, offset int) ([]models.Product, int64, error)
	UpdateProduct(ctx context.Context, p *models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uint) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(r repositories.ProductRepository) ProductService {
	return &productService{repo: r}
}

func (s *productService) CreateProduct(ctx context.Context, p *models.Product) (*models.Product, error) {
	if p.PriceCents <= 0 {
		return nil, ErrInvalidPrice
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productService) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *productService) ListProducts(ctx context.Context, limit, offset int) ([]models.Product, int64, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *productService) UpdateProduct(ctx context.Context, p *models.Product) (*models.Product, error) {
	if p.PriceCents <= 0 {
		return nil, ErrInvalidPrice
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
