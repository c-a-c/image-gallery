// @title: common.go
// @author: @Asameshi00
// @date: 2025-10-10
// @description: 画像ドメインモデル

package domain

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

// Image
// @description: ユーザーがアップロードした画像の型定義
type Image struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null"`
	User         User           `json:"user" gorm:"foreignKey:UserID"`
	Title        string         `json:"title" gorm:"not null"`
	Description  string         `json:"description"`
	CloudinaryID string         `json:"cloudinary_id" gorm:"not null"`
	URL          string         `json:"url" gorm:"not null"`
	Width        int            `json:"width"`
	Height       int            `json:"height"`
	FileSize     int64          `json:"file_size"`
	Format       string         `json:"format"`
	Tags         string         `json:"tags"` // Comma-separated tags
	IsPublic     bool           `json:"is_public" gorm:"default:true"`
	ViewCount    int            `json:"view_count" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Post
// @description: 投稿を含む画像
type Post struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Images      []Image        `json:"images" gorm:"many2many:post_images;"`
	Tags        string         `json:"tags"` // Comma-separated tags
	IsPublic    bool           `json:"is_public" gorm:"default:true"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ImageRepository
// @description: 画像データ操作のインターフェース
type ImageRepository interface {
	Create(image *Image) error // 画像を作成
	GetByID(id uint) (*Image, error) // 画像をIDで取得
	GetByUserID(userID uint, offset, limit int) ([]*Image, error) // ユーザーIDで画像を取得
	GetPublic(offset, limit int) ([]*Image, error) // 公開画像を取得
	Update(image *Image) error // 画像を更新
	Delete(id uint) error // 画像を削除
	Search(query string, offset, limit int) ([]*Image, error) // 画像をクエリで検索
	GetByTags(tags []string, offset, limit int) ([]*Image, error) // タグで画像を取得
	IncrementViewCount(id uint) error // 閲覧数を増やす
}

// PostRepository
// @description: 投稿データ操作のインターフェース
type PostRepository interface {
	Create(post *Post) error // 投稿を作成
	GetByID(id uint) (*Post, error) // 投稿をIDで取得
	GetByUserID(userID uint, offset, limit int) ([]*Post, error) // ユーザーIDで投稿を取得
	GetPublic(offset, limit int) ([]*Post, error) // 公開投稿を取得
	Update(post *Post) error // 投稿を更新
	Delete(id uint) error // 投稿を削除
	Search(query string, offset, limit int) ([]*Post, error) // 投稿をクエリで検索
	GetByTags(tags []string, offset, limit int) ([]*Post, error) // タグで投稿を取得
	IncrementViewCount(id uint) error // 閲覧数を増やす
}

// ImageUseCase
// @description: 画像ビジネスロジックのインターフェース
type ImageUseCase interface {
	UploadImage(userID uint, title, description, tags string, imageData []byte, filename string) (*Image, error) // 画像をアップロード
	UploadImageFromFile(userID uint, title, description, tags string, file *multipart.FileHeader) (*Image, error) // 画像をファイルからアップロード
	GetImage(imageID uint) (*Image, error) // 画像をIDで取得
	GetUserImages(userID uint, page, limit int) ([]*Image, error) // ユーザーIDで画像を取得
	GetPublicImages(page, limit int) ([]*Image, error) // 公開画像を取得
	UpdateImage(userID, imageID uint, title, description, tags string, isPublic bool) (*Image, error) // 画像を更新
	DeleteImage(userID, imageID uint) error // 画像を削除
	SearchImages(query string, page, limit int) ([]*Image, error) // 画像をクエリで検索
	GetImagesByTags(tags []string, page, limit int) ([]*Image, error) // タグで画像を取得
	IncrementViewCount(imageID uint) error // 閲覧数を増やす
}

// PostUseCase
// @description: 投稿ビジネスロジックのインターフェース
type PostUseCase interface {
	CreatePost(userID uint, title, description string, imageIDs []uint, tags string) (*Post, error) // 投稿を作成
	GetPost(postID uint) (*Post, error) // 投稿をIDで取得
	GetUserPosts(userID uint, page, limit int) ([]*Post, error) // ユーザーIDで投稿を取得
	GetPublicPosts(page, limit int) ([]*Post, error) // 公開投稿を取得
	UpdatePost(userID, postID uint, title, description, tags string, isPublic bool) (*Post, error) // 投稿を更新
	DeletePost(userID, postID uint) error // 投稿を削除
	SearchPosts(query string, page, limit int) ([]*Post, error) // 投稿をクエリで検索
	GetPostsByTags(tags []string, page, limit int) ([]*Post, error) // タグで投稿を取得
	IncrementViewCount(postID uint) error // 閲覧数を増やす
}
