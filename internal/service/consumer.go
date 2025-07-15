package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"github.com/yoanesber/go-consumer-api-with-jwt/config/database"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/repository"
)

// Interface for consumer service
// This interface defines the methods that the consumer service should implement
type ConsumerService interface {
	GetAllConsumers(page int, limit int) ([]entity.Consumer, error)
	GetConsumerByID(id string) (entity.Consumer, error)
	GetActiveConsumers(page int, limit int) ([]entity.Consumer, error)
	GetInactiveConsumers(page int, limit int) ([]entity.Consumer, error)
	GetSuspendedConsumers(page int, limit int) ([]entity.Consumer, error)
	CreateConsumer(c entity.Consumer) (entity.Consumer, error)
	UpdateConsumerStatus(id string, status string) (entity.Consumer, error)
}

// This struct defines the ConsumerService that contains a repository field of type ConsumerRepository
// It implements the ConsumerService interface and provides methods for consumer-related operations
type consumerService struct {
	repo repository.ConsumerRepository
}

// NewConsumerService creates a new instance of ConsumerService with the given repository.
// This function initializes the consumerService struct and returns it.
func NewConsumerService(repo repository.ConsumerRepository) ConsumerService {
	return &consumerService{repo: repo}
}

// GetAllConsumers retrieves all consumers from the database.
func (s *consumerService) GetAllConsumers(page int, limit int) ([]entity.Consumer, error) {
	db := database.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Retrieve all consumers from the repository
	consumers, err := s.repo.GetAllConsumers(db, page, limit)
	if err != nil {
		return nil, err
	}

	return consumers, nil
}

// GetConsumerByID retrieves a consumer by its ID from the database.
func (s *consumerService) GetConsumerByID(id string) (entity.Consumer, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.Consumer{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the consumer by ID from the repository
	consumer, err := s.repo.GetConsumerByID(db, id)
	if err != nil {
		return entity.Consumer{}, err
	}

	return consumer, nil
}

// GetActiveConsumers retrieves all active consumers from the database.
func (s *consumerService) GetActiveConsumers(page int, limit int) ([]entity.Consumer, error) {
	db := database.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Retrieve all active consumers from the repository
	activeConsumers, err := s.repo.GetConsumersByStatus(db, entity.ConsumerStatusActive, page, limit)
	if err != nil {
		return nil, err
	}

	return activeConsumers, nil
}

// GetInactiveConsumers retrieves all inactive consumers from the database.
func (s *consumerService) GetInactiveConsumers(page int, limit int) ([]entity.Consumer, error) {
	db := database.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Retrieve all inactive consumers from the repository
	inactiveConsumers, err := s.repo.GetConsumersByStatus(db, "inactive", page, limit)
	if err != nil {
		return nil, err
	}

	return inactiveConsumers, nil
}

// GetSuspendedConsumers retrieves all suspended consumers from the database.
func (s *consumerService) GetSuspendedConsumers(page int, limit int) ([]entity.Consumer, error) {
	db := database.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Retrieve all suspended consumers from the repository
	suspendedConsumers, err := s.repo.GetConsumersByStatus(db, "suspended", page, limit)
	if err != nil {
		return nil, err
	}

	return suspendedConsumers, nil
}

// CreateConsumer creates a new consumer in the database.
// It validates the consumer struct and checks if the ID already exists before creating a new consumer.
func (s *consumerService) CreateConsumer(c entity.Consumer) (entity.Consumer, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.Consumer{}, fmt.Errorf("database connection is nil")
	}

	// Validate the consumer struct using the validator
	if err := c.Validate(); err != nil {
		return entity.Consumer{}, err
	}

	createdConsumer := entity.Consumer{}
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the username already exists
		existingConsumer, err := s.repo.GetConsumerByUsername(db, c.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check existing consumer by username: %w", err)
		}

		// If the consumer already exists, return an error
		if (err == nil) || !(existingConsumer.Equals(&entity.Consumer{})) {
			return fmt.Errorf("consumer with username %s already exists", c.Username)
		}

		// Check if the email already exists
		existingConsumer, err = s.repo.GetConsumerByEmail(db, c.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check existing consumer by email: %w", err)
		}

		// If the consumer already exists, return an error
		if (err == nil) || !(existingConsumer.Equals(&entity.Consumer{})) {
			return fmt.Errorf("consumer with email %s already exists", c.Email)
		}

		// Check if the phone already exists
		normalizedPhone := NormalizePhoneNumber(c.Phone)
		c.Phone = normalizedPhone
		existingConsumer, err = s.repo.GetConsumerByPhone(db, normalizedPhone)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check existing consumer by phone: %w", err)
		}

		// If the consumer already exists, return an error
		if (err == nil) || !(existingConsumer.Equals(&entity.Consumer{})) {
			return fmt.Errorf("consumer with phone %s already exists", c.Phone)
		}

		c.Status = "inactive" // Set default status to inactive
		createdConsumer, err = s.repo.CreateConsumer(tx, c)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Consumer{}, err
	}

	return createdConsumer, nil
}

// NormalizePhoneNumber removes non-digit characters and ensures country code (e.g., starts with 62)
func NormalizePhoneNumber(phone string) string {
	// Remove all non-digit characters
	re := regexp.MustCompile(`\D`)
	digitsOnly := re.ReplaceAllString(phone, "")

	// Replace leading '0' with '62' (Indonesia)
	if strings.HasPrefix(digitsOnly, "0") {
		digitsOnly = "62" + digitsOnly[1:]
	}

	return digitsOnly
}

// UpdateConsumerStatus updates the status of an existing consumer in the database.
// It checks if the consumer exists and validates the status before updating it.
func (s *consumerService) UpdateConsumerStatus(id string, status string) (entity.Consumer, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.Consumer{}, fmt.Errorf("database connection is nil")
	}

	updatedConsumer := entity.Consumer{}
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the consumer exists
		existingConsumer, err := s.repo.GetConsumerByID(db, id)
		if err != nil {
			return err
		}

		existingConsumer.Status = status
		updatedConsumer, err = s.repo.UpdateConsumer(tx, existingConsumer)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Consumer{}, err
	}

	return updatedConsumer, nil
}
