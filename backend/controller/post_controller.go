package controller

import (
	"backend/domain"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// PostController handles post requests
type PostController struct {
	postUseCase domain.PostUseCase
}

// NewPostController creates a new post controller
func NewPostController(postUseCase domain.PostUseCase) *PostController {
	return &PostController{
		postUseCase: postUseCase,
	}
}

// CreatePost handles post creation
func (c *PostController) CreatePost(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	var req domain.PostCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	post, err := c.postUseCase.CreatePost(userID, req.Title, req.Description, req.ImageIDs, req.Tags)
	if err != nil {
		switch err {
		case domain.ErrForbidden:
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": "You don't have permission to use one or more of the specified images",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create post",
			})
		}
	}

	return ctx.JSON(http.StatusCreated, post)
}

// GetPost handles getting a single post
func (c *PostController) GetPost(ctx echo.Context) error {
	postIDStr := ctx.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid post ID",
		})
	}

	post, err := c.postUseCase.GetPost(uint(postID))
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Post not found",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get post",
			})
		}
	}

	return ctx.JSON(http.StatusOK, post)
}

// GetUserPosts handles getting user's posts
func (c *PostController) GetUserPosts(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	page, limit := getPaginationParams(ctx)
	posts, err := c.postUseCase.GetUserPosts(userID, page, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get posts",
		})
	}

	return ctx.JSON(http.StatusOK, posts)
}

// GetPublicPosts handles getting public posts
func (c *PostController) GetPublicPosts(ctx echo.Context) error {
	page, limit := getPaginationParams(ctx)
	posts, err := c.postUseCase.GetPublicPosts(page, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get posts",
		})
	}

	return ctx.JSON(http.StatusOK, posts)
}

// UpdatePost handles updating a post
func (c *PostController) UpdatePost(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	postIDStr := ctx.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid post ID",
		})
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Tags        string `json:"tags"`
		IsPublic    bool   `json:"is_public"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	post, err := c.postUseCase.UpdatePost(userID, uint(postID), req.Title, req.Description, req.Tags, req.IsPublic)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Post not found",
			})
		case domain.ErrForbidden:
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": "You don't have permission to update this post",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update post",
			})
		}
	}

	return ctx.JSON(http.StatusOK, post)
}

// DeletePost handles deleting a post
func (c *PostController) DeletePost(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	postIDStr := ctx.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid post ID",
		})
	}

	err = c.postUseCase.DeletePost(userID, uint(postID))
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Post not found",
			})
		case domain.ErrForbidden:
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": "You don't have permission to delete this post",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete post",
			})
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Post deleted successfully",
	})
}

// SearchPosts handles searching posts
func (c *PostController) SearchPosts(ctx echo.Context) error {
	query := ctx.QueryParam("q")
	if query == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required",
		})
	}

	page, limit := getPaginationParams(ctx)
	posts, err := c.postUseCase.SearchPosts(query, page, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to search posts",
		})
	}

	return ctx.JSON(http.StatusOK, posts)
}

// GetPostsByTags handles getting posts by tags
func (c *PostController) GetPostsByTags(ctx echo.Context) error {
	tagsParam := ctx.QueryParam("tags")
	if tagsParam == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tags parameter is required",
		})
	}

	tags := strings.Split(tagsParam, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	page, limit := getPaginationParams(ctx)
	posts, err := c.postUseCase.GetPostsByTags(tags, page, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get posts by tags",
		})
	}

	return ctx.JSON(http.StatusOK, posts)
}
