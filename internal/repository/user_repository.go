package repositories

import (
	"context"

	"github.com/AlexRijikov/go-petshop-api/internal/models"
	"gorm.io/gorm"
)

// UserRepository визначає методи для доступу до даних користувача
// Використовує GORM для взаємодії з базою даних
// Всі методи використовують WithContext(ctx) — корисно для таймаутів/тестів.

// UserRepository визначає методи для роботи з користувачами (створення, пошук за email).

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error                    // створює нового користувача
	GetByEmail(ctx context.Context, email string) (*models.User, error)  // шукає користувача за email
	GetByUsername(username string) (*models.User, error)                 // шукає користувача за username
	GetByID(id uint) (*models.User, error)                               // шукає користувача за ID
	UpdateProfile(id uint, username, email string) (*models.User, error) // оновлює дані користувача (username, email)
}

// userRepo реалізує UserRepository

type userRepo struct {
	db *gorm.DB
}

// NewUserRepository створює новий UserRepository

func NewUserRepository(db *gorm.DB) UserRepository { // приймає підключення до бази даних
	return &userRepo{db: db} // повертаємо новий екземпляр userRepo з підключенням до бази даних
}

// Create додає нового користувача в базу даних

func (r *userRepo) Create(ctx context.Context, u *models.User) error { // приймає контекст і користувача для створення
	return r.db.WithContext(ctx).Create(u).Error // створюємо нового користувача в базі даних
}

// GetByEmail шукає користувача за email і повертає його або помилку, якщо не знайдено

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User // створюємо змінну для збереження знайденого користувача
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err // якщо користувача не знайдено, повертаємо помилку
	}
	return &user, nil // повертаємо знайденого користувача
}

// GetByUsername шукає користувача за username і повертає його або помилку, якщо не знайдено

func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID шукає користувача за ID і повертає його або помилку, якщо не знайдено
func (r *userRepo) GetByID(id uint) (*models.User, error) {

	var user models.User                                // створюємо змінну для збереження знайденого користувача
	if err := r.db.First(&user, id).Error; err != nil { // шукаємо користувача за ID
		return nil, err // якщо користувача не знайдено, повертаємо помилку
	}
	return &user, nil // повертаємо знайденого користувача
}

// UpdateProfile оновлює дані користувача (username, email) за його ID
func (r *userRepo) UpdateProfile(id uint, username, email string) (*models.User, error) {
	var user models.User                                // створюємо змінну для збереження користувача
	if err := r.db.First(&user, id).Error; err != nil { // шукаємо користувача за ID
		return nil, err // якщо користувача не знайдено, повертаємо помилку
	}

	// Оновлюємо поля користувача

	user.Username = username
	user.Email = email

	// Зберігаємо оновленого користувача в базу даних

	if err := r.db.Save(&user).Error; err != nil { // зберігаємо оновленого користувача
		return nil, err // якщо сталася помилка при збереженні, повертаємо її
	}
	return &user, nil // повертаємо оновленого користувача

}
