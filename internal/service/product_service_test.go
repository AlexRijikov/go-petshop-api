package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"github.com/AlexRijikov/go-petshop-api/internal/service"
)

// Тестуємо бізнес-логіку без DB.

// Простий in-memory repo реалізує repositories.ProductRepository для тесту

type memRepo struct {
	data map[uint]*models.Product
	next uint
}

func newMemRepo() *memRepo {
	return &memRepo{data: map[uint]*models.Product{}, next: 1}
}

func (m *memRepo) Create(ctx context.Context, p *models.Product) error {
	p.ID = m.next
	m.next++
	m.data[p.ID] = p
	return nil
}
func (m *memRepo) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	p, ok := m.data[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}
func (m *memRepo) List(ctx context.Context, limit, offset int) ([]models.Product, int64, error) {
	var out []models.Product
	for _, v := range m.data {
		out = append(out, *v)
	}
	return out, int64(len(out)), nil
}
func (m *memRepo) Update(ctx context.Context, p *models.Product) error {
	m.data[p.ID] = p
	return nil
}
func (m *memRepo) Delete(ctx context.Context, id uint) error {
	delete(m.data, id)
	return nil
}

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
