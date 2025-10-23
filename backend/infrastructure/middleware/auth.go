package middleware

import (
	"backend/usecase"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header is required",
				})
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Token is required",
				})
			}

			// Validate token
			claims, err := usecase.ValidateJWTToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_username", claims.Username)

			return next(c)
		}
	}
}

// OptionalAuthMiddleware handles optional JWT authentication
func OptionalAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return next(c)
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return next(c)
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return next(c)
			}

			// Validate token
			claims, err := usecase.ValidateJWTToken(token)
			if err != nil {
				return next(c)
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_username", claims.Username)

			return next(c)
		}
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusOK)
			}

			return next(c)
		}
	}
}

// ErrorHandler handles errors
func ErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		c.JSON(he.Code, map[string]string{
			"error": he.Message.(string),
		})
		return
	}

	c.JSON(http.StatusInternalServerError, map[string]string{
		"error": "Internal server error",
	})
}

// LoggerMiddleware logs requests
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Log request
			c.Logger().Infof("Request: %s %s", c.Request().Method, c.Request().URL.Path)

			err := next(c)

			// Log response
			c.Logger().Infof("Response: %d", c.Response().Status)

			return err
		}
	}
}
