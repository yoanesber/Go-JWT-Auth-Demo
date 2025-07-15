package test_consumer

import (
	"time"

	"gorm.io/gorm" // Import GORM for ORM functionalities

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
)

// ConsumerMockedRepository is an interface that defines the methods for interacting with consumer data in a mocked repository.
// It includes methods for retrieving, creating, and updating consumers in the database.
type ConsumerMockedRepository interface {
	GetAllConsumers(tx *gorm.DB, page int, limit int) ([]entity.Consumer, error)
	GetConsumerByID(tx *gorm.DB, id string) (entity.Consumer, error)
	GetConsumerByUsername(tx *gorm.DB, username string) (entity.Consumer, error)
	GetConsumerByEmail(tx *gorm.DB, email string) (entity.Consumer, error)
	GetConsumerByPhone(tx *gorm.DB, phone string) (entity.Consumer, error)
	GetConsumersByStatus(tx *gorm.DB, status string, page int, limit int) ([]entity.Consumer, error)
	CreateConsumer(tx *gorm.DB, d entity.Consumer) (entity.Consumer, error)
	UpdateConsumer(tx *gorm.DB, d entity.Consumer) (entity.Consumer, error)
}

// consumerMockedRepository is a struct that implements the ConsumerMockedRepository interface.
// It contains methods for interacting with consumer data in a mocked repository.
type consumerMockedRepository struct{}

// NewConsumerMockedRepository creates a new instance of ConsumerMockedRepository.
// It initializes the consumerMockedRepository struct and returns it.
func NewConsumerMockedRepository() ConsumerMockedRepository {
	return &consumerMockedRepository{}
}

// GetAllConsumers retrieves all consumers from the dummy data.
// It simulates the retrieval of consumer data from a database by returning a predefined list of consumers
func (r *consumerMockedRepository) GetAllConsumers(tx *gorm.DB, page int, limit int) ([]entity.Consumer, error) {
	return getDummyConsumers(), nil
}

// GetConsumerByID retrieves a consumer by its ID from the dummy data.
// It simulates the retrieval of a single consumer from a database by returning a predefined consumer object
func (r *consumerMockedRepository) GetConsumerByID(tx *gorm.DB, id string) (entity.Consumer, error) {
	if id == "" {
		return entity.Consumer{}, nil // Return an empty consumer if ID is empty
	}

	consumer := getDummyConsumer()
	if consumer.ID != id {
		return entity.Consumer{}, gorm.ErrRecordNotFound // Return an error if the ID does
	}

	return consumer, nil
}

// GetConsumerByUsername retrieves a consumer by its username from the dummy data.
// It simulates the retrieval of a single consumer from a database by returning a predefined consumer object
func (r *consumerMockedRepository) GetConsumerByUsername(tx *gorm.DB, username string) (entity.Consumer, error) {
	if username == "" {
		return entity.Consumer{}, nil // Return an empty consumer if username is empty
	}

	consumer := getDummyConsumer()
	if consumer.Username != username {
		return entity.Consumer{}, gorm.ErrRecordNotFound // Return an error if the username does not match
	}

	return consumer, nil
}

// GetConsumerByEmail retrieves a consumer by its email from the dummy data.
// It simulates the retrieval of a single consumer from a database by returning a predefined consumer object
func (r *consumerMockedRepository) GetConsumerByEmail(tx *gorm.DB, email string) (entity.Consumer, error) {
	if email == "" {
		return entity.Consumer{}, nil // Return an empty consumer if email is empty
	}

	consumer := getDummyConsumer()
	if consumer.Email != email {
		return entity.Consumer{}, gorm.ErrRecordNotFound // Return an error if the email does not match
	}

	return consumer, nil
}

// GetConsumerByPhone retrieves a consumer by its phone number from the dummy data.
// It simulates the retrieval of a single consumer from a database by returning a predefined consumer object
func (r *consumerMockedRepository) GetConsumerByPhone(tx *gorm.DB, phone string) (entity.Consumer, error) {
	if phone == "" {
		return entity.Consumer{}, nil // Return an empty consumer if phone is empty
	}

	consumer := getDummyConsumer()
	if consumer.Phone != phone {
		return entity.Consumer{}, gorm.ErrRecordNotFound // Return an error if the phone does not match
	}

	return consumer, nil
}

// GetConsumersByStatus retrieves consumers by their status from the dummy data.
// It simulates the retrieval of a list of consumers from a database by filtering the predefined list
func (r *consumerMockedRepository) GetConsumersByStatus(tx *gorm.DB, status string, page int, limit int) ([]entity.Consumer, error) {
	consumers := getDummyConsumers()
	var filteredConsumers []entity.Consumer

	for _, consumer := range consumers {
		if consumer.Status == status {
			filteredConsumers = append(filteredConsumers, consumer)
		}
	}

	return filteredConsumers, nil
}

// CreateConsumer creates a new consumer in the dummy data.
// It simulates the creation of a consumer in a database by returning a predefined consumer object
func (r *consumerMockedRepository) CreateConsumer(tx *gorm.DB, t entity.Consumer) (entity.Consumer, error) {
	if t.ID == "" {
		t.ID = "new-dummy-id" // Assign a new ID if not provided
	}
	t.CreatedAt = time.Now()
	t.UpdatedAt = t.CreatedAt

	return t, nil
}

// UpdateConsumer updates an existing consumer in the dummy data.
// It simulates the update of a consumer in a database by returning a predefined consumer object
func (r *consumerMockedRepository) UpdateConsumer(tx *gorm.DB, t entity.Consumer) (entity.Consumer, error) {
	consumer := getDummyConsumer()
	consumer.ID = t.ID
	consumer.Fullname = t.Fullname
	consumer.Username = t.Username
	consumer.Email = t.Email
	consumer.Phone = t.Phone
	consumer.Address = t.Address
	consumer.BirthDate = t.BirthDate
	consumer.Status = t.Status
	consumer.UpdatedAt = time.Now()

	return consumer, nil
}
