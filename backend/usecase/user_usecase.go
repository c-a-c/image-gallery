package usecase

import (
	"backend/domain"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// userUseCase
// @description: ユーザーユースケースの実装
type userUseCase struct {
	userRepo domain.UserRepository
}

// NewUserUseCase
// @description: ユーザーユースケースを初期化
func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

// Register
// @description: 新規ユーザーを登録
func (u *userUseCase) Register(email, username, password, firstName, lastName string) (*domain.User, error) {
	// メールアドレスがすでに存在するかどうかを確認
	existingUser, err := u.userRepo.GetByEmail(email)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrEmailExists
	}

	// ユーザー名がすでに存在するかどうかを確認
	existingUser, err = u.userRepo.GetByUsername(username)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrUsernameExists
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// ユーザーを作成
	user := &domain.User{
		Email:     email,
		Username:  username,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	err = u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login
// @description: ユーザーを認証
func (u *userUseCase) Login(email, password string) (*domain.User, string, error) {
	// メールアドレスでユーザーを取得
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", domain.ErrInvalidPassword
	}

	// ユーザーがアクティブかどうかを確認
	if !user.IsActive {
		return nil, "", domain.ErrForbidden
	}

	// パスワードを確認
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", domain.ErrInvalidPassword
	}

	// JWTトークンを生成
	token, err := u.generateJWTToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetProfile
// @description: ユーザープロフィールを取得
func (u *userUseCase) GetProfile(userID uint) (*domain.User, error) {
	return u.userRepo.GetByID(userID)
}

// UpdateProfile
// @description: ユーザープロフィールを更新
func (u *userUseCase) UpdateProfile(userID uint, firstName, lastName, avatar string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Avatar = avatar

	err = u.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword
// @description: ユーザーパスワードを変更
func (u *userUseCase) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// 古いパスワードを確認
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return domain.ErrInvalidPassword
	}

	// 新しいパスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	return u.userRepo.Update(user)
}

// DeactivateAccount
// @description: ユーザーアカウントを非アクティブ化
func (u *userUseCase) DeactivateAccount(userID uint) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.IsActive = false
	return u.userRepo.Update(user)
}

// generateJWTToken
// @description: ユーザーのJWTトークンを生成
func (u *userUseCase) generateJWTToken(user *domain.User) (string, error) {
	claims := &domain.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Exp:      time.Now().Add(time.Hour * 24).Unix(),
		Iat:      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // 本番環境では強力なシークレットキーを使用する
	}

	return token.SignedString([]byte(jwtSecret))
}

// ValidateJWTToken validates a JWT token
func ValidateJWTToken(tokenString string) (*domain.JWTClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // In production, this should be a strong secret
	}

	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, domain.ErrUnauthorized
}

// GenerateJWTToken generates a JWT token for the user
func GenerateJWTToken(user *domain.User) (string, error) {
	claims := &domain.JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Exp:      time.Now().Add(time.Hour * 24).Unix(),
		Iat:      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // In production, this should be a strong secret
	}

	return token.SignedString([]byte(jwtSecret))
}
