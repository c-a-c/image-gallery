package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // Password is excluded from JSON
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Avatar    string         `json:"avatar"` // Cloudinary URL
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(offset, limit int) ([]*User, error)
}

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	Register(email, username, password, firstName, lastName string) (*User, error)
	Login(email, password string) (*User, string, error) // Returns user and JWT token
	GetProfile(userID uint) (*User, error)
	UpdateProfile(userID uint, firstName, lastName, avatar string) (*User, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
	DeactivateAccount(userID uint) error
}
