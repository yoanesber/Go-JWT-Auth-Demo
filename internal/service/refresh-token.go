package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/yoanesber/go-consumer-api-with-jwt/config/database"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/repository"
)

// Interface for refresh token service
// This interface defines the methods that the refresh token service should implement
type RefreshTokenService interface {
	GetRefreshTokenByUserID(userID int64) (entity.RefreshToken, error)
	GetRefreshTokenByToken(token string) (entity.RefreshToken, error)
	VerifyExpirationDate(exp time.Time) (bool, error)
	CreateRefreshToken(userID int64) (entity.RefreshToken, error)
}

// This struct defines the RefreshTokenService that contains a repository field of type RefreshTokenRepository
// It implements the RefreshTokenService interface and provides methods for refresh token-related operations
type refreshTokenService struct {
	repo repository.RefreshTokenRepository
}

// NewRefreshTokenService creates a new instance of RefreshTokenService with the given repository.
// It initializes the refreshTokenService struct and returns it.
func NewRefreshTokenService(repo repository.RefreshTokenRepository) RefreshTokenService {
	return &refreshTokenService{repo: repo}
}

// GetRefreshTokenByUserID retrieves a refresh token by its user ID from the database.
func (s *refreshTokenService) GetRefreshTokenByUserID(userID int64) (entity.RefreshToken, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.RefreshToken{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the token by user ID from the repository
	token, err := s.repo.GetRefreshTokenByUserID(db, userID)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return token, nil
}

// GetRefreshTokenByToken retrieves a refresh token by its token string from the database.
func (s *refreshTokenService) GetRefreshTokenByToken(token string) (entity.RefreshToken, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.RefreshToken{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the token by token string from the repository
	refreshToken, err := s.repo.GetRefreshTokenByToken(db, token)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return refreshToken, nil
}

// VerifyExpirationDate checks if the expiration date is valid and not in the past.
func (s *refreshTokenService) VerifyExpirationDate(exp time.Time) (bool, error) {
	// Check if the expiration date is valid
	if exp.IsZero() {
		return false, fmt.Errorf("expiration date is not set")
	}

	// Check if the expiration date is in the past
	if time.Now().After(exp) {
		return false, nil
	}

	return true, nil
}

// CreateRefreshToken creates a new refresh token for the user in the database.
// If a refresh token already exists for the user, it will be removed before creating a new one,
// ensuring that only one refresh token exists for each user at a time.
func (s *refreshTokenService) CreateRefreshToken(userID int64) (entity.RefreshToken, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.RefreshToken{}, fmt.Errorf("database connection is nil")
	}

	createdRefreshToken := entity.RefreshToken{}
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the refresh token already exists for the user
		existingRefreshToken, err := s.repo.GetRefreshTokenByUserID(tx, userID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// If the refresh token already exists, remove it
		if !existingRefreshToken.Equals(&entity.RefreshToken{}) {
			if _, err := s.repo.RemoveRefreshTokenByUserID(tx, userID); err != nil {
				return err
			}
		}

		// Create a new refresh token
		tokenStr := uuid.New().String()
		refreshToken := entity.RefreshToken{
			Token:      tokenStr,
			UserID:     userID,
			ExpiryDate: GetRefreshTokenExpiration(time.Now()),
		}

		// Create the refresh token in the database
		createdRefreshToken, err = s.repo.CreateRefreshToken(tx, refreshToken)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.RefreshToken{}, err
	}

	return createdRefreshToken, nil
}

// GetRefreshTokenExpiration calculates the expiration date for the refresh token.
// It retrieves the expiration hour from an environment variable and adds it to the current time.
func GetRefreshTokenExpiration(now time.Time) time.Time {
	expHour, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION_HOUR"))
	if err != nil {
		return now.Add(24 * time.Hour) // Default to 24 hours if the environment variable is not set or invalid
	}
	if expHour <= 0 {
		expHour = 24
	}

	return now.Add(time.Hour * time.Duration(expHour))
}
