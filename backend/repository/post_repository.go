package repository

import (
	"backend/domain"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// postRepository implements domain.PostRepository
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) domain.PostRepository {
	return &postRepository{db: db}
}

// Create creates a new post
func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

// GetByID retrieves a post by ID
func (r *postRepository) GetByID(id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.db.Preload("User").Preload("Images").First(&post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &post, nil
}

// GetByUserID retrieves posts by user ID
func (r *postRepository) GetByUserID(userID uint, offset, limit int) ([]*domain.Post, error) {
	var posts []*domain.Post
	err := r.db.Where("user_id = ?", userID).
		Preload("Images").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// GetPublic retrieves public posts
func (r *postRepository) GetPublic(offset, limit int) ([]*domain.Post, error) {
	var posts []*domain.Post
	err := r.db.Where("is_public = ?", true).
		Preload("User").
		Preload("Images").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// Update updates a post
func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

// Delete deletes a post
func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Post{}, id).Error
}

// Search searches posts by query
func (r *postRepository) Search(query string, offset, limit int) ([]*domain.Post, error) {
	var posts []*domain.Post
	searchQuery := "%" + strings.ToLower(query) + "%"

	err := r.db.Where("is_public = ? AND (LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(tags) LIKE ?)",
		true, searchQuery, searchQuery, searchQuery).
		Preload("User").
		Preload("Images").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// GetByTags retrieves posts by tags
func (r *postRepository) GetByTags(tags []string, offset, limit int) ([]*domain.Post, error) {
	var posts []*domain.Post

	query := r.db.Where("is_public = ?", true)

	for _, tag := range tags {
		query = query.Where("LOWER(tags) LIKE ?", "%"+strings.ToLower(tag)+"%")
	}

	err := query.Preload("User").
		Preload("Images").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// IncrementViewCount increments the view count for a post
func (r *postRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&domain.Post{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}
