package test_consumer

import (
	"time"

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/customtype"
)

// getDummyConsumer returns a dummy consumer entity for testing purposes.
func getDummyConsumer() entity.Consumer {
	return entity.Consumer{
		ID:        "dummy-id",
		Fullname:  "Dummy Consumer",
		Username:  "dummyuser",
		Email:     "dummy-user@example.com",
		Phone:     "6281234567890",
		Address:   "123 Dummy Street",
		BirthDate: &customtype.Date{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// getDummyConsumers returns a slice of dummy consumer entities for testing purposes.
func getDummyConsumers() []entity.Consumer {
	return []entity.Consumer{
		{
			ID:        "dummy-id-1",
			Fullname:  "Dummy Consumer 1",
			Username:  "dummyuser1",
			Email:     "dummy-user-1@example.com",
			Phone:     "6281234567891",
			Address:   "123 Dummy Street 1",
			BirthDate: &customtype.Date{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "dummy-id-2",
			Fullname:  "Dummy Consumer 2",
			Username:  "dummyuser2",
			Email:     "dummy-user-2@example.com",
			Phone:     "6281234567892",
			Address:   "123 Dummy Street 2",
			BirthDate: &customtype.Date{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},
			Status:    "inactive",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "dummy-id-3",
			Fullname:  "Dummy Consumer 3",
			Username:  "dummyuser3",
			Email:     "dummy-user-3@example.com",
			Phone:     "6281234567893",
			Address:   "123 Dummy Street 3",
			BirthDate: &customtype.Date{Time: time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC)},
			Status:    "suspended",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "dummy-id-4",
			Fullname:  "Dummy Consumer 4",
			Username:  "dummyuser4",
			Email:     "dummy-user-4@example.com",
			Phone:     "6281234567894",
			Address:   "123 Dummy Street 4",
			BirthDate: &customtype.Date{Time: time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC)},
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "dummy-id-5",
			Fullname:  "Dummy Consumer 5",
			Username:  "dummyuser5",
			Email:     "dummy-user-5@example.com",
			Phone:     "6281234567895",
			Address:   "123 Dummy Street 5",
			BirthDate: &customtype.Date{Time: time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC)},
			Status:    "inactive",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
