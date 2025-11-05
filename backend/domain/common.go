// @title: common.go
// @author: @Asameshi00
// @date: 2025-10-10
// @description: 普段使うエラーのドメインモデル

package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// @description: 普段使うエラーの定義
var (
	ErrNotFound        = errors.New("resource not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidInput    = errors.New("invalid input")
	ErrEmailExists     = errors.New("email already exists")
	ErrUsernameExists  = errors.New("username already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrFileTooLarge    = errors.New("file too large")
	ErrInvalidFileType = errors.New("invalid file type")
	ErrUploadFailed    = errors.New("upload failed")
)

// @description: ページネーションリクエスト
type PaginationRequest struct {
	Page  int `json:"page" form:"page" query:"page"`
	Limit int `json:"limit" form:"limit" query:"limit"`
}

// @description: ページネーションレスポンス
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// @description: 認証リクエスト
type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// @description: ユーザー登録リクエスト
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=20"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// @description: 認証レスポンス
type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// @description: 画像アップロードリクエスト
type ImageUploadRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	IsPublic    bool   `json:"is_public"`
}

// @description: 投稿作成リクエスト
type PostCreateRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	ImageIDs    []uint `json:"image_ids" validate:"required,min=1"`
	Tags        string `json:"tags"`
	IsPublic    bool   `json:"is_public"`
}

// @description: 検索リクエスト
type SearchRequest struct {
	Query string `json:"query" form:"query" query:"query"`
	Page  int    `json:"page" form:"page" query:"page"`
	Limit int    `json:"limit" form:"limit" query:"limit"`
}

// @description: アプリケーション設定
type Config struct {
	Port                string
	DBHost              string
	DBPort              string
	DBName              string
	DBUser              string
	DBPassword          string
	DBMaxOpenConns      int
	DBMaxIdleConns      int
	DBConnMaxLifetime   int
	JWTSecret           string
	CloudinaryURL       string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	CloudinaryCloudName string
}

// @description: JWTトークンクレーム
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

// @description: jwt.Claims.GetAudienceを実装
func (c *JWTClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

// @description: jwt.Claims.GetExpirationTimeを実装
func (c *JWTClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.Exp, 0)), nil
}

// @description: jwt.Claims.GetIssuedAtを実装
func (c *JWTClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.Iat, 0)), nil
}

// @description: jwt.Claims.GetIssuerを実装
func (c *JWTClaims) GetIssuer() (string, error) {
	return "", nil
}

// @description: jwt.Claims.GetNotBeforeを実装
func (c *JWTClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// @description: jwt.Claims.GetSubjectを実装
func (c *JWTClaims) GetSubject() (string, error) {
	return "", nil
}

// @description: トークンが有効かどうかを確認
func (c *JWTClaims) Valid() error {
	if time.Now().Unix() > c.Exp {
		return ErrUnauthorized
	}
	return nil
}
