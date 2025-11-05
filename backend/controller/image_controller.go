// @title: image_controller.go
// @author: @Asameshi00
// @date: 2025-10-10
// @description: 画像コントローラー

package controller

import (
	"backend/domain"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ImageController handles image requests
type ImageController struct {
	imageUseCase domain.ImageUseCase
}

// NewImageController creates a new image controller
func NewImageController(imageUseCase domain.ImageUseCase) *ImageController {
	return &ImageController{
		imageUseCase: imageUseCase,
	}
}

// UploadImage handles image upload
func (c *ImageController) UploadImage(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	// 画像ファイルを取得
	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "No image file provided",
		})
	}

	// フォームデータを取得
	title := ctx.FormValue("title")
	description := ctx.FormValue("description")
	tags := ctx.FormValue("tags")

	image, err := c.imageUseCase.UploadImageFromFile(userID, title, description, tags, file)
	if err != nil {
		switch err {
		case domain.ErrInvalidFileType:
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid file type. Only images are allowed",
			})
		case domain.ErrFileTooLarge:
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "File too large. Maximum size is 10MB",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to upload image",
			})
		}
	}

	return ctx.JSON(http.StatusCreated, image)
}

// GetImage handles getting a single image
func (c *ImageController) GetImage(ctx echo.Context) error {
	imageIDStr := ctx.Param("id")
	imageID, err := strconv.ParseUint(imageIDStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid image ID",
		})
	}

	image, err := c.imageUseCase.GetImage(uint(imageID))
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Image not found",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get image",
			})
		}
	}

	return ctx.JSON(http.StatusOK, image)
}

// GetUserImages handles getting user's images
func (c *ImageController) GetUserImages(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	page, limit := getPaginationParams(ctx)
	images, err := c.imageUseCase.GetUserImages(userID, page, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get images",
		})
	}

	return ctx.JSON(http.StatusOK, images)
}

// GetPublicImages handles getting public images
func (c *ImageController) GetPublicImages(ctx echo.Context) error {
	page, limit := getPaginationParams(ctx)
	images, err := c.imageUseCase.GetPublicImages(page, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get images",
		})
	}

	return ctx.JSON(http.StatusOK, images)
}

// UpdateImage handles updating an image
func (c *ImageController) UpdateImage(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	imageIDStr := ctx.Param("id")
	imageID, err := strconv.ParseUint(imageIDStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid image ID",
		})
	}

	var req domain.ImageUploadRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	image, err := c.imageUseCase.UpdateImage(userID, uint(imageID), req.Title, req.Description, req.Tags, req.IsPublic)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Image not found",
			})
		case domain.ErrForbidden:
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": "You don't have permission to update this image",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update image",
			})
		}
	}

	return ctx.JSON(http.StatusOK, image)
}

// DeleteImage handles deleting an image
func (c *ImageController) DeleteImage(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	imageIDStr := ctx.Param("id")
	imageID, err := strconv.ParseUint(imageIDStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid image ID",
		})
	}

	err = c.imageUseCase.DeleteImage(userID, uint(imageID))
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Image not found",
			})
		case domain.ErrForbidden:
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": "You don't have permission to delete this image",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete image",
			})
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Image deleted successfully",
	})
}

// SearchImages handles searching images
func (c *ImageController) SearchImages(ctx echo.Context) error {
	query := ctx.QueryParam("q")
	if query == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required",
		})
	}

	page, limit := getPaginationParams(ctx)
	images, err := c.imageUseCase.SearchImages(query, page, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to search images",
		})
	}

	return ctx.JSON(http.StatusOK, images)
}

// getPaginationParams extracts pagination parameters from request
func getPaginationParams(ctx echo.Context) (int, int) {
	pageStr := ctx.QueryParam("page")
	limitStr := ctx.QueryParam("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return page, limit
}
