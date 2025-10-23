package repository

import (
	"backend/domain"
	"errors"

	"gorm.io/gorm"
)

// userRepository implements domain.UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

// List retrieves users with pagination
func (r *userRepository) List(offset, limit int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// CheckEmailExists checks if email already exists
func (r *userRepository) CheckEmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// CheckUsernameExists checks if username already exists
func (r *userRepository) CheckUsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
