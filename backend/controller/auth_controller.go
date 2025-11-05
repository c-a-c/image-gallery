package controller

import (
	"backend/domain"
	"backend/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthController handles authentication requests
type AuthController struct {
	userUseCase domain.UserUseCase
}

// NewAuthController creates a new auth controller
func NewAuthController(userUseCase domain.UserUseCase) *AuthController {
	return &AuthController{
		userUseCase: userUseCase,
	}
}

// Register handles user registration
func (c *AuthController) Register(ctx echo.Context) error {
	var req domain.RegisterRequest
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

	user, err := c.userUseCase.Register(req.Email, req.Username, req.Password, req.FirstName, req.LastName)
	if err != nil {
		switch err {
		case domain.ErrEmailExists, domain.ErrUsernameExists:
			return ctx.JSON(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to register user",
			})
		}
	}

	// Generate JWT token
	token, err := usecase.GenerateJWTToken(user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	return ctx.JSON(http.StatusCreated, domain.AuthResponse{
		User:  user,
		Token: token,
	})
}

// Login handles user login
func (c *AuthController) Login(ctx echo.Context) error {
	var req domain.AuthRequest
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

	user, token, err := c.userUseCase.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidPassword, domain.ErrForbidden:
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid credentials",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to login",
			})
		}
	}

	return ctx.JSON(http.StatusOK, domain.AuthResponse{
		User:  user,
		Token: token,
	})
}

// GetProfile handles getting user profile
func (c *AuthController) GetProfile(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	user, err := c.userUseCase.GetProfile(userID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return ctx.JSON(http.StatusOK, user)
}

// UpdateProfile handles updating user profile
func (c *AuthController) UpdateProfile(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	user, err := c.userUseCase.UpdateProfile(userID, req.FirstName, req.LastName, req.Avatar)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update profile",
		})
	}

	return ctx.JSON(http.StatusOK, user)
}

// ChangePassword handles changing user password
func (c *AuthController) ChangePassword(ctx echo.Context) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	err = c.userUseCase.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case domain.ErrInvalidPassword:
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid old password",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to change password",
			})
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}

// getUserIDFromContext extracts user ID from JWT token in context
func getUserIDFromContext(ctx echo.Context) (uint, error) {
	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return 0, domain.ErrUnauthorized
	}
	return userID, nil
}
