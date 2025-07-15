package entity

import (
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/customtype"
)

const (
	ConsumerStatusActive    = "active"
	ConsumerStatusInactive  = "inactive"
	ConsumerStatusSuspended = "suspended"
)

// Consumer represents the consumer entity in the database.
type Consumer struct {
	ID        string           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Fullname  string           `gorm:"type:varchar(100);not null" json:"fullname" validate:"required,max=100"`
	Username  string           `gorm:"type:varchar(50);unique;not null" json:"username" validate:"required,max=50"`
	Email     string           `gorm:"type:varchar(100);unique;not null" json:"email" validate:"required,email,max=100"`
	Phone     string           `gorm:"type:varchar(20);unique;not null" json:"phone" validate:"required,max=20"`
	Address   string           `gorm:"type:text;not null" json:"address" validate:"required"`
	BirthDate *customtype.Date `gorm:"type:date" json:"birthDate,omitempty" validate:"required,omitempty"`
	Status    string           `gorm:"type:varchar(20);not null;default:'inactive';check:status IN ('active','inactive','suspended')" json:"status"`
	CreatedAt time.Time        `gorm:"column:created_at;type:timestamptz;autoCreateTime;default:now()" json:"createdAt,omitempty"`
	UpdatedAt time.Time        `gorm:"column:updated_at;type:timestamptz;autoUpdateTime;default:now()" json:"updatedAt,omitempty"`
}

// TableName overrides the table name used by Consumer to `consumers`.
func (Consumer) TableName() string {
	return "consumers"
}

// Equals compares two Consumer objects for equality.
func (c *Consumer) Equals(other *Consumer) bool {
	if c == nil && other == nil {
		return true
	}

	if c == nil || other == nil {
		return false
	}

	if (c.ID != other.ID) ||
		(c.Fullname != other.Fullname) ||
		(c.Username != other.Username) ||
		(c.Email != other.Email) ||
		(c.Phone != other.Phone) ||
		(c.Address != other.Address) ||
		(c.BirthDate != nil && other.BirthDate != nil && *c.BirthDate != *other.BirthDate) ||
		(c.Status != other.Status) {
		return false
	}

	return true
}

// Validate validates the Consumer struct using the validator package.
func (c *Consumer) Validate() error {
	var v *validator.Validate = validator.New()

	if err := v.Struct(c); err != nil {
		return err
	}
	return nil
}
