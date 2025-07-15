package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
)

// Interface for user repository
// This interface defines the methods that the user repository should implement
type UserRepository interface {
	GetUserByID(tx *gorm.DB, id int64) (entity.User, error)
	GetUserByUsername(tx *gorm.DB, username string) (entity.User, error)
	GetUserByEmail(tx *gorm.DB, email string) (entity.User, error)
	UpdateUser(tx *gorm.DB, user entity.User) (entity.User, error)
}

// This struct defines the UserRepository that contains methods for interacting with the database
// It implements the UserRepository interface and provides methods for user-related operations
type userRepository struct{}

// NewUserRepository creates a new instance of UserRepository.
// It initializes the userRepository struct and returns it.
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// GetUserByID retrieves a user by its ID from the database.
func (r *userRepository) GetUserByID(tx *gorm.DB, id int64) (entity.User, error) {
	// Select the user with the given ID from the database
	var user entity.User
	err := tx.Preload("Roles").First(&user, "id = ?", id).Error

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username from the database.
func (r *userRepository) GetUserByUsername(tx *gorm.DB, username string) (entity.User, error) {
	// Select the user with the given username from the database
	var user entity.User
	err := tx.Preload("Roles").First(&user, "lower(username) = lower(?)", username).Error

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email from the database.
func (r *userRepository) GetUserByEmail(tx *gorm.DB, email string) (entity.User, error) {
	// Select the user with the given email from the database
	var user entity.User
	err := tx.Preload("Roles").First(&user, "lower(email) = lower(?)", email).Error

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// UpdateUser updates an existing user in the database and returns the updated user.
func (r *userRepository) UpdateUser(tx *gorm.DB, user entity.User) (entity.User, error) {
	// Update the user in the database
	if err := tx.Save(&user).Error; err != nil {
		return entity.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
