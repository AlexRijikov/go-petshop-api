package services
import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"github.com/AlexRijikov/go-petshop-api/internal/repository"
)

// CartService визначає інтерфейс для роботи з кошиком користувача

type CartService interface {
	List(ctx context.Context, userID uint) ([]models.CartItem, error)
	AddItem(ctx context.Context, userID, productID uint, quantity int) error
	UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error
	RemoveItem(ctx context.Context, userID, productID uint) error
}
// cartService реалізує інтерфейс CartService

type cartService struct {
	cartRepo repositories.CartRepository 
}

// NewCartService створює новий CartService з наданим репозиторієм	

func NewCartService(cr repositories.CartRepository) CartService {
	return &cartService{cartRepo: cr}
}

// List отримує всі товари в кошику користувача

func (s *cartService) List(ctx context.Context, userID uint) ([]models.CartItem, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID") 
	}
	items, err := s.cartRepo.ListByUserID(ctx, userID) 
	if err != nil {
		log.Printf("Error fetching cart items for user %d: %v", userID, err) 
		return nil, fmt.Errorf("failed to fetch cart items: %w", err)
	} 
	return items, nil
}

// AddItem -  

func (s *cartService) AddItem(ctx context.Context, userID, productID uint, quantity int) error {
	if userID == 0 || productID == 0 || quantity <= 0 {
		return errors.New("invalid input data")
	}
	item := &models.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.cartRepo.Add(ctx, item)
}

// UpdateQuantity 

func (s *cartService) UpdateQuantity(ctx context.Context, userID, productID uint, quantity int) error {
	if userID == 0 || productID == 0 || quantity <= 0 {
		return errors.New("invalid input data")
	}
	return s.cartRepo.UpdateQuantity(ctx, userID, productID, quantity)
}
 
//RemoveItem - 

func (s *cartService) RemoveItem(ctx context.Context, userID, productID uint) error {
	if userID == 0 || productID == 0 {
		return errors.New("invalid input data")
	}
	return s.cartRepo.Remove(ctx, userID, productID)
}