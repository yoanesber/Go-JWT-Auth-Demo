package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
)

// Interface for refresh token repository
// This interface defines the methods that the refresh token repository should implement
type RefreshTokenRepository interface {
	GetRefreshTokenByUserID(tx *gorm.DB, userID int64) (entity.RefreshToken, error)
	GetRefreshTokenByToken(tx *gorm.DB, token string) (entity.RefreshToken, error)
	CreateRefreshToken(tx *gorm.DB, token entity.RefreshToken) (entity.RefreshToken, error)
	RemoveRefreshTokenByUserID(tx *gorm.DB, userID int64) (bool, error)
}

// This struct defines the RefreshTokenRepository that contains methods for interacting with the database
// It implements the RefreshTokenRepository interface and provides methods for refresh token-related operations
type refreshTokenRepository struct{}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository.
// It initializes the refreshTokenRepository struct and returns it.
func NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepository{}
}

// GetRefreshTokenByUserID retrieves a refresh token by its user ID from the database.
func (r *refreshTokenRepository) GetRefreshTokenByUserID(tx *gorm.DB, userID int64) (entity.RefreshToken, error) {
	// Select the refresh token with the given user ID from the database
	var refreshToken entity.RefreshToken
	err := tx.First(&refreshToken, "user_id = ?", userID).Error
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return refreshToken, nil
}

// GetRefreshTokenByToken retrieves a refresh token by its token string from the database.
func (r *refreshTokenRepository) GetRefreshTokenByToken(tx *gorm.DB, token string) (entity.RefreshToken, error) {
	// Select the refresh token with the given token string from the database
	var refreshToken entity.RefreshToken
	err := tx.First(&refreshToken, "token = ?", token).Error
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return refreshToken, nil
}

// CreateRefreshToken creates a new refresh token in the database.
func (r *refreshTokenRepository) CreateRefreshToken(tx *gorm.DB, token entity.RefreshToken) (entity.RefreshToken, error) {
	// Create a new refresh token in the database
	if err := tx.Create(&token).Error; err != nil {
		return entity.RefreshToken{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return token, nil
}

// RemoveRefreshTokenByUserID removes a refresh token by its user ID from the database.
func (r *refreshTokenRepository) RemoveRefreshTokenByUserID(tx *gorm.DB, userID int64) (bool, error) {
	// Delete the refresh token with the given user ID from the database
	if err := tx.Where("user_id = ?", userID).Delete(&entity.RefreshToken{}).Error; err != nil {
		return false, fmt.Errorf("failed to remove refresh token by user ID %d: %w", userID, err)
	}

	return true, nil
}
