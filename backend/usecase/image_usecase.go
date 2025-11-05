// @title: image_usecase.go
// @author: @Asameshi00
// @date: 2025-10-23
// @description: 画像ユースケースの実装

package usecase

import (
	"slices"
	"backend/domain"
	"backend/infrastructure/cloudinary"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// imageUseCase
// @description: 画像ユースケースの実装
type imageUseCase struct {
	imageRepo     domain.ImageRepository
	cloudinarySvc *cloudinary.Service
}

// NewImageUseCase
// @description: 画像ユースケースを初期化
func NewImageUseCase(imageRepo domain.ImageRepository) domain.ImageUseCase {
	cloudinarySvc, err := cloudinary.NewService()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Cloudinary service: %v", err))
	}

	return &imageUseCase{
		imageRepo:     imageRepo,
		cloudinarySvc: cloudinarySvc,
	}
}

// UploadImage
// @description: 画像をアップロード
func (u *imageUseCase) UploadImage(userID uint, title, description, tags string, imageData []byte, filename string) (*domain.Image, error) {
	// Validate file type
	if !isValidImageType(filename) {
		return nil, domain.ErrInvalidFileType
	}

	// 画像サイズは100MBまで
	if len(imageData) > 100*1024*1024 {
		return nil, domain.ErrFileTooLarge
	}

	// Cloudinaryにアップロード
	ctx := context.Background()
	result, err := u.cloudinarySvc.UploadImage(ctx, strings.NewReader(string(imageData)), filename, "images")
	if err != nil {
		return nil, fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	// 画像レコードを作成
	image := &domain.Image{
		UserID:       userID,
		Title:        title,
		Description:  description,
		Tags:         tags,
		CloudinaryID: result.PublicID,
		URL:          result.SecureURL,
		Width:        result.Width,
		Height:       result.Height,
		FileSize:     int64(result.Bytes),
		Format:       result.Format,
		IsPublic:     true,
		ViewCount:    0,
	}

	err = u.imageRepo.Create(image)
	if err != nil {
		// データベース保存に失敗した場合、Cloudinaryから削除
		u.cloudinarySvc.DeleteImage(ctx, result.PublicID)
		return nil, err
	}

	return image, nil
}

// GetImage
// @description: 画像をIDで取得
func (u *imageUseCase) GetImage(imageID uint) (*domain.Image, error) {
	image, err := u.imageRepo.GetByID(imageID)
	if err != nil {
		return nil, err
	}

	// 閲覧数を増やす
	go u.imageRepo.IncrementViewCount(imageID)

	return image, nil
}

// GetUserImages
// @description: ユーザーIDで画像を取得
func (u *imageUseCase) GetUserImages(userID uint, page, limit int) ([]*domain.Image, error) {
	offset := (page - 1) * limit
	return u.imageRepo.GetByUserID(userID, offset, limit)
}

// GetPublicImages
// @description: 公開画像を取得
func (u *imageUseCase) GetPublicImages(page, limit int) ([]*domain.Image, error) {
	offset := (page - 1) * limit
	return u.imageRepo.GetPublic(offset, limit)
}

// UpdateImage
// @description: 画像を更新
func (u *imageUseCase) UpdateImage(
	userID,
	imageID uint,
	title,
	description,
	tags string,
	isPublic bool,
) (*domain.Image, error) {
	image, err := u.imageRepo.GetByID(imageID)
	if err != nil {
		return nil, err
	}

	// ユーザーが画像の所有者かどうかを確認
	if image.UserID != userID {
		return nil, domain.ErrForbidden
	}

	image.Title = title
	image.Description = description
	image.Tags = tags
	image.IsPublic = isPublic

	err = u.imageRepo.Update(image)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// DeleteImage
// @description: 画像を削除
func (u *imageUseCase) DeleteImage(userID, imageID uint) error {
	image, err := u.imageRepo.GetByID(imageID)
	if err != nil {
		return err
	}

	// ユーザーが画像の所有者かどうかを確認
	if image.UserID != userID {
		return domain.ErrForbidden
	}

	// Cloudinaryから削除
	ctx := context.Background()
	err = u.cloudinarySvc.DeleteImage(ctx, image.CloudinaryID)
	if err != nil {
		// ログを出力してデータベース削除を続行
		fmt.Printf("Failed to delete image from Cloudinary: %v\n", err)
	}

	// データベースから削除
	return u.imageRepo.Delete(imageID)
}

// SearchImages
// @description: 画像をクエリで検索
func (u *imageUseCase) SearchImages(query string, page, limit int) ([]*domain.Image, error) {
	offset := (page - 1) * limit
	return u.imageRepo.Search(query, offset, limit)
}

// GetImagesByTags
// @description: タグで画像を取得
func (u *imageUseCase) GetImagesByTags(tags []string, page, limit int) ([]*domain.Image, error) {
	offset := (page - 1) * limit
	return u.imageRepo.GetByTags(tags, offset, limit)
}

// IncrementViewCount
// @description: 閲覧数を増やす
func (u *imageUseCase) IncrementViewCount(imageID uint) error {
	return u.imageRepo.IncrementViewCount(imageID)
}

// isValidImageType
// @description: ファイルタイプが有効かどうかを確認
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}

	return slices.Contains(validExts, ext)
}

// UploadImageFromFile
// @description: マルチパートファイルから画像をアップロード
func (u *imageUseCase) UploadImageFromFile(userID uint, title, description, tags string, file *multipart.FileHeader) (*domain.Image, error) {
	// ファイルを開く
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// ファイルデータを読み込む
	fileData := make([]byte, file.Size)
	_, err = src.Read(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return u.UploadImage(userID, title, description, tags, fileData, file.Filename)
}
