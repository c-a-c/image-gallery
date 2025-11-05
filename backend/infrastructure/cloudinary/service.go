// @title: service.go
// @author: @Asameshi00
// @date: 2025-10-23
// @description: Cloudinaryサービスの実装

package cloudinary

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Service
// @description: Cloudinaryサービスの実装
type Service struct {
	cld *cloudinary.Cloudinary
}

// NewService
// @description: Cloudinaryサービスを初期化
func NewService() (*Service, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("missing Cloudinary credentials")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &Service{cld: cld}, nil
}

// UploadImage
// @description: 画像をCloudinaryにアップロード
func (s *Service) UploadImage(ctx context.Context, imageData io.Reader, filename string, folder string) (*uploader.UploadResult, error) {
	uploadParams := uploader.UploadParams{
		Folder:       folder,
		PublicID:     filename,
		ResourceType: "image",
	}

	result, err := s.cld.Upload.Upload(ctx, imageData, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	return result, nil
}

// DeleteImage
// @description: 画像をCloudinaryから削除
func (s *Service) DeleteImage(ctx context.Context, publicID string) error {
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

// GetImageURL
// @description: 画像のCloudinaryURLを生成
func (s *Service) GetImageURL(publicID string, transformations map[string]interface{}) string {
	img, err := s.cld.Image(publicID)
	if err != nil {
		return ""
	}
	return img.Config.Cloud.APIKey
}

// TransformImage
// @description: 画像を変換したURLを生成
func (s *Service) TransformImage(publicID string, width, height int, crop string) string {
	img, err := s.cld.Image(publicID)
	if err != nil {
		return ""
	}

	return img.AssetType.String()
}
