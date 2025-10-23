package usecase

import (
	"backend/domain"
	"fmt"
	"strings"
)

// postUseCase
// @description: 投稿ユースケースの実装
type postUseCase struct {
	postRepo  domain.PostRepository
	imageRepo domain.ImageRepository
}

// NewPostUseCase
// @description: 投稿ユースケースを初期化
func NewPostUseCase(postRepo domain.PostRepository, imageRepo domain.ImageRepository) domain.PostUseCase {
	return &postUseCase{
		postRepo:  postRepo,
		imageRepo: imageRepo,
	}
}

// CreatePost
// @description: 投稿を作成
func (u *postUseCase) CreatePost(userID uint, title, description string, imageIDs []uint, tags string) (*domain.Post, error) {
	// すべての画像がユーザーのものかどうかを確認
	for _, imageID := range imageIDs {
		image, err := u.imageRepo.GetByID(imageID)
		if err != nil {
			return nil, fmt.Errorf("image not found: %d", imageID)
		}
		if image.UserID != userID {
			return nil, domain.ErrForbidden
		}
	}

	// 投稿を作成
	post := &domain.Post{
		UserID:      userID,
		Title:       title,
		Description: description,
		Tags:        tags,
		IsPublic:    true,
		ViewCount:   0,
	}

	err := u.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	// Associate images with post
	// Note: This would require updating the post_images many-to-many relationship
	// For now, we'll just return the post without the images association
	// In a real implementation, you'd need to handle the many-to-many relationship

	return post, nil
}

// GetPost
// @description: 投稿をIDで取得
func (u *postUseCase) GetPost(postID uint) (*domain.Post, error) {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	// 閲覧数を増やす
	go u.postRepo.IncrementViewCount(postID)

	return post, nil
}

// GetUserPosts
// @description: ユーザーIDで投稿を取得
func (u *postUseCase) GetUserPosts(userID uint, page, limit int) ([]*domain.Post, error) {
	offset := (page - 1) * limit
	return u.postRepo.GetByUserID(userID, offset, limit)
}

// GetPublicPosts
// @description: 公開投稿を取得
func (u *postUseCase) GetPublicPosts(page, limit int) ([]*domain.Post, error) {
	offset := (page - 1) * limit
	return u.postRepo.GetPublic(offset, limit)
}

// UpdatePost
// @description: 投稿を更新
func (u *postUseCase) UpdatePost(userID, postID uint, title, description, tags string, isPublic bool) (*domain.Post, error) {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	// ユーザーが投稿の所有者かどうかを確認
	if post.UserID != userID {
		return nil, domain.ErrForbidden
	}

	post.Title = title
	post.Description = description
	post.Tags = tags
	post.IsPublic = isPublic

	err = u.postRepo.Update(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// DeletePost
// @description: 投稿を削除
func (u *postUseCase) DeletePost(userID, postID uint) error {
	post, err := u.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	// ユーザーが投稿の所有者かどうかを確認
	if post.UserID != userID {
		return domain.ErrForbidden
	}

	return u.postRepo.Delete(postID)
}

// SearchPosts
// @description: 投稿をクエリで検索
func (u *postUseCase) SearchPosts(query string, page, limit int) ([]*domain.Post, error) {
	offset := (page - 1) * limit
	return u.postRepo.Search(query, offset, limit)
}

// GetPostsByTags
// @description: タグで投稿を取得
func (u *postUseCase) GetPostsByTags(tags []string, page, limit int) ([]*domain.Post, error) {
	offset := (page - 1) * limit
	return u.postRepo.GetByTags(tags, offset, limit)
}

// IncrementViewCount
// @description: 閲覧数を増やす
func (u *postUseCase) IncrementViewCount(postID uint) error {
	return u.postRepo.IncrementViewCount(postID)
}

// parseTags
// @description: カンマ区切りのタグをスライスにパース
func parseTags(tags string) []string {
	if tags == "" {
		return []string{}
	}

	tagSlice := strings.Split(tags, ",")
	var result []string
	for _, tag := range tagSlice {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
