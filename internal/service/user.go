package service

import (
	"fmt"
	"time"

	"github.com/yoanesber/go-consumer-api-with-jwt/config/database"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/repository"
	"gorm.io/gorm"
)

// Interface for user service
// This interface defines the methods that the user service should implement
type UserService interface {
	GetUserByID(id int64) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	UpdateLastLogin(id int64, lastLogin time.Time) (bool, error)
}

// This struct defines the UserService that contains a repository field of type UserRepository
// It implements the UserService interface and provides methods for user-related operations
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService with the given repository.
// It initializes the userService struct and returns it.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// GetUserByID retrieves a user by its ID from the database.
func (s *userService) GetUserByID(id int64) (entity.User, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.User{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the user by ID from the repository
	user, err := s.repo.GetUserByID(db, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username from the database.
func (s *userService) GetUserByUsername(username string) (entity.User, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.User{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the user by username from the repository
	user, err := s.repo.GetUserByUsername(db, username)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email from the database.
func (s *userService) GetUserByEmail(email string) (entity.User, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.User{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the user by email from the repository
	user, err := s.repo.GetUserByEmail(db, email)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// UpdateLastLogin updates the last login time of a user in the database.
func (s *userService) UpdateLastLogin(id int64, lastLogin time.Time) (bool, error) {
	db := database.GetPostgres()
	if db == nil {
		return false, fmt.Errorf("database connection is nil")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the user exists
		existingUser, err := s.repo.GetUserByID(db, id)
		if err != nil {
			return err
		}

		// Check if the existing user is empty
		if (existingUser.Equals(&entity.User{})) {
			return fmt.Errorf("user with ID %d not found", id)
		}

		// Update the last login time
		*existingUser.LastLogin = lastLogin
		_, err = s.repo.UpdateUser(tx, existingUser)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
