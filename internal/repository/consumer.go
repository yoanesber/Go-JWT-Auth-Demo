package repository

import (
	"fmt"

	"gorm.io/gorm" // Import GORM for ORM functionalities

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
)

// Interface for consumer repository
// This interface defines the methods that the consumer repository should implement
type ConsumerRepository interface {
	GetAllConsumers(tx *gorm.DB, page int, limit int) ([]entity.Consumer, error)
	GetConsumerByID(tx *gorm.DB, id string) (entity.Consumer, error)
	GetConsumerByUsername(tx *gorm.DB, username string) (entity.Consumer, error)
	GetConsumerByEmail(tx *gorm.DB, email string) (entity.Consumer, error)
	GetConsumerByPhone(tx *gorm.DB, phone string) (entity.Consumer, error)
	GetConsumersByStatus(tx *gorm.DB, status string, page int, limit int) ([]entity.Consumer, error)
	CreateConsumer(tx *gorm.DB, d entity.Consumer) (entity.Consumer, error)
	UpdateConsumer(tx *gorm.DB, d entity.Consumer) (entity.Consumer, error)
}

// This struct defines the consumerRepository that implements the ConsumerRepository interface.
// It contains methods for interacting with the consumer data in the database.
type consumerRepository struct{}

// NewConsumerRepository creates a new instance of ConsumerRepository.
// It initializes the consumerRepository struct and returns it.
func NewConsumerRepository() ConsumerRepository {
	return &consumerRepository{}
}

// GetAllConsumers retrieves all consumers from the database.
func (r *consumerRepository) GetAllConsumers(tx *gorm.DB, page int, limit int) ([]entity.Consumer, error) {
	var consumers []entity.Consumer
	err := tx.Order("created_at ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&consumers).Error

	if err != nil {
		return nil, err
	}

	return consumers, nil
}

// It returns a single consumer by its ID from the database.
func (r *consumerRepository) GetConsumerByID(tx *gorm.DB, id string) (entity.Consumer, error) {
	var consumer entity.Consumer
	err := tx.First(&consumer, "id = ?", id).Error

	if err != nil {
		return entity.Consumer{}, err
	}

	return consumer, nil
}

// GetConsumerByEmail retrieves a consumer by their email from the database.
func (r *consumerRepository) GetConsumerByUsername(tx *gorm.DB, username string) (entity.Consumer, error) {
	var consumer entity.Consumer
	err := tx.First(&consumer, "lower(username) = lower(?)", username).Error

	if err != nil {
		return entity.Consumer{}, err
	}

	return consumer, nil
}

// GetConsumerByEmail retrieves a consumer by their email from the database.
func (r *consumerRepository) GetConsumerByEmail(tx *gorm.DB, email string) (entity.Consumer, error) {
	var consumer entity.Consumer
	err := tx.First(&consumer, "lower(email) = lower(?)", email).Error

	if err != nil {
		return entity.Consumer{}, err
	}

	return consumer, nil
}

// GetConsumerByPhone retrieves a consumer by their phone number from the database.
func (r *consumerRepository) GetConsumerByPhone(tx *gorm.DB, phone string) (entity.Consumer, error) {
	var consumer entity.Consumer
	err := tx.First(&consumer, "phone = ?", phone).Error

	if err != nil {
		return entity.Consumer{}, err
	}

	return consumer, nil
}

// GetActiveConsumers retrieves all active consumers from the database.
func (r *consumerRepository) GetConsumersByStatus(tx *gorm.DB, status string, page int, limit int) ([]entity.Consumer, error) {
	var consumers []entity.Consumer
	err := tx.Where("status = ?", status).
		Order("created_at ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&consumers).
		Error

	if err != nil {
		return nil, err
	}

	return consumers, nil
}

// CreateConsumer creates a new consumer in the database and returns the created consumer.
func (r *consumerRepository) CreateConsumer(tx *gorm.DB, t entity.Consumer) (entity.Consumer, error) {
	// Insert new consumer
	if err := tx.Create(&t).Error; err != nil {
		return entity.Consumer{}, fmt.Errorf("failed to create consumer: %w", err)
	}

	return t, nil
}

// UpdateConsumer updates an existing consumer in the database and returns the updated consumer.
// This method is used to modify an existing consumer's details.
func (r *consumerRepository) UpdateConsumer(tx *gorm.DB, t entity.Consumer) (entity.Consumer, error) {
	// Save the updated consumer
	if err := tx.Save(&t).Error; err != nil {
		return entity.Consumer{}, fmt.Errorf("failed to update consumer: %w", err)
	}

	return t, nil
}
