package repository

import (
	"backend/domain"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// imageRepository implements domain.ImageRepository
type imageRepository struct {
	db *gorm.DB
}

// NewImageRepository creates a new image repository
func NewImageRepository(db *gorm.DB) domain.ImageRepository {
	return &imageRepository{db: db}
}

// Create creates a new image
func (r *imageRepository) Create(image *domain.Image) error {
	return r.db.Create(image).Error
}

// GetByID retrieves an image by ID
func (r *imageRepository) GetByID(id uint) (*domain.Image, error) {
	var image domain.Image
	err := r.db.Preload("User").First(&image, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &image, nil
}

// GetByUserID retrieves images by user ID
func (r *imageRepository) GetByUserID(userID uint, offset, limit int) ([]*domain.Image, error) {
	var images []*domain.Image
	err := r.db.Where("user_id = ?", userID).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&images).Error
	return images, err
}

// GetPublic retrieves public images
func (r *imageRepository) GetPublic(offset, limit int) ([]*domain.Image, error) {
	var images []*domain.Image
	err := r.db.Where("is_public = ?", true).
		Preload("User").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&images).Error
	return images, err
}

// Update updates an image
func (r *imageRepository) Update(image *domain.Image) error {
	return r.db.Save(image).Error
}

// Delete deletes an image
func (r *imageRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Image{}, id).Error
}

// Search searches images by query
func (r *imageRepository) Search(query string, offset, limit int) ([]*domain.Image, error) {
	var images []*domain.Image
	searchQuery := "%" + strings.ToLower(query) + "%"

	err := r.db.Where("is_public = ? AND (LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(tags) LIKE ?)",
		true, searchQuery, searchQuery, searchQuery).
		Preload("User").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&images).Error
	return images, err
}

// GetByTags retrieves images by tags
func (r *imageRepository) GetByTags(tags []string, offset, limit int) ([]*domain.Image, error) {
	var images []*domain.Image

	query := r.db.Where("is_public = ?", true)

	for _, tag := range tags {
		query = query.Where("LOWER(tags) LIKE ?", "%"+strings.ToLower(tag)+"%")
	}

	err := query.Preload("User").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&images).Error
	return images, err
}

// IncrementViewCount increments the view count for an image
func (r *imageRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&domain.Image{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}
