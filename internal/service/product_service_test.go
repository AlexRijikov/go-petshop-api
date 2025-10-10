package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"github.com/AlexRijikov/go-petshop-api/internal/service"
)

// Тестуємо бізнес-логіку без БД, використовуючи in-memory репозиторій.

// Простий in-memory repo реалізує repositories.ProductRepository для тестів.

type memRepo struct {
	data map[uint]*models.Product
	next uint
}

// newMemRepo створює новий in-memory репозиторій

func newMemRepo() *memRepo {
	return &memRepo{data: map[uint]*models.Product{}, next: 1}
}

// Реалізація методів ProductRepository для memRepo

func (m *memRepo) Create(ctx context.Context, p *models.Product) error {
	p.ID = m.next
	m.next++
	m.data[p.ID] = p
	return nil
}

// GetByID повертає продукт за ID або nil, nil якщо не знайдено

func (m *memRepo) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	p, ok := m.data[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

// List повертає всі продукти (без реальної пагінації для простоти)

func (m *memRepo) List(ctx context.Context, limit, offset int) ([]models.Product, int64, error) {
	var out []models.Product
	for _, v := range m.data {
		out = append(out, *v)
	}
	return out, int64(len(out)), nil
}

// Update оновлює продукт в пам'яті (перезаписує по ID)

func (m *memRepo) Update(ctx context.Context, p *models.Product) error {
	m.data[p.ID] = p
	return nil
}

// Delete видаляє продукт з пам'яті

func (m *memRepo) Delete(ctx context.Context, id uint) error {
	delete(m.data, id)
	return nil
}

// Тести для ProductService

func TestCreateProduct(t *testing.T) {
	repo := newMemRepo()
	svc := services.NewProductService(repo)

	p := &models.Product{
		Name:       "Test",
		PriceCents: 100,
		Stock:      1,
	}
	created, err := svc.CreateProduct(context.Background(), p)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, uint(1), created.ID)
}

// Тест створення продукту з некоректною ціною

func TestCreateInvalidPrice(t *testing.T) {
	repo := newMemRepo()
	svc := services.NewProductService(repo)

	p := &models.Product{
		Name:       "Bad",
		PriceCents: 0,
		Stock:      1,
	}
	_, err := svc.CreateProduct(context.Background(), p)
	assert.Error(t, err)
}
