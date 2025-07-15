package entity

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"

	validation "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/validation-util"
)

// User represents the user entity in the database.
type User struct {
	ID                        int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username                  string          `gorm:"type:varchar(20);not null;unique" json:"username" validate:"required,min=3,max=20"`
	Password                  string          `gorm:"type:varchar(150);not null" json:"password" validate:"required,min=8"`
	Email                     string          `gorm:"type:varchar(100);not null;unique" json:"email" validate:"required,email,max=100"`
	Firstname                 string          `gorm:"type:varchar(20);not null" json:"firstName" validate:"required,max=20"`
	Lastname                  *string         `gorm:"type:varchar(20)" json:"lastName,omitempty" validate:"omitempty,max=20"`
	IsEnabled                 *bool           `gorm:"not null;default:false" json:"isEnabled,omitempty"`
	IsAccountNonExpired       *bool           `gorm:"not null;default:false" json:"isAccountNonExpired,omitempty"`
	IsAccountNonLocked        *bool           `gorm:"not null;default:false" json:"isAccountNonLocked,omitempty"`
	IsCredentialsNonExpired   *bool           `gorm:"not null;default:false" json:"isCredentialsNonExpired,omitempty"`
	IsDeleted                 *bool           `gorm:"not null;default:false" json:"isDeleted,omitempty"`
	AccountExpirationDate     *time.Time      `gorm:"type:timestamptz" json:"accountExpirationDate,omitempty"`
	CredentialsExpirationDate *time.Time      `gorm:"type:timestamptz" json:"credentialsExpirationDate,omitempty"`
	UserType                  string          `gorm:"type:varchar(20);not null;check:user_type IN ('SERVICE_ACCOUNT','USER_ACCOUNT')" json:"userType" validate:"required,max=20,oneof=SERVICE_ACCOUNT USER_ACCOUNT"`
	LastLogin                 *time.Time      `json:"lastLogin,omitempty"`
	CreatedBy                 *int64          `json:"createdBy,omitempty"`
	CreatedAt                 *time.Time      `gorm:"type:timestamptz;autoCreateTime;default:now()" json:"createdAt,omitempty"`
	UpdatedBy                 *int64          `json:"updatedBy,omitempty"`
	UpdatedAt                 *time.Time      `gorm:"type:timestamptz;autoUpdateTime;default:now()" json:"updatedAt,omitempty"`
	DeletedBy                 *int64          `json:"deletedBy,omitempty"`
	DeletedAt                 *gorm.DeletedAt `gorm:"type:timestamptz;index" json:"deletedAt,omitempty"`
	Roles                     []Role          `gorm:"many2many:user_roles;constraint:OnUpdate:RESTRICT,OnDelete:SET NULL" json:"roles,omitempty"`
}

// Override the TableName method to specify the table name
// in the database. This is optional if you want to use the default naming convention.
func (User) TableName() string {
	return "users"
}

// Equals compares two User objects for equality.
func (u *User) Equals(other *User) bool {
	if u == nil && other == nil {
		return true
	}

	if u == nil || other == nil {
		return false
	}

	if (u.ID != other.ID) ||
		(u.Username != other.Username) ||
		(u.Password != other.Password) ||
		(u.Email != other.Email) ||
		(u.Firstname != other.Firstname) ||
		(u.Lastname != other.Lastname) ||
		(u.IsEnabled != other.IsEnabled) ||
		(u.IsAccountNonExpired != other.IsAccountNonExpired) ||
		(u.IsAccountNonLocked != other.IsAccountNonLocked) ||
		(u.IsCredentialsNonExpired != other.IsCredentialsNonExpired) ||
		(u.IsDeleted != other.IsDeleted) ||
		(u.AccountExpirationDate != other.AccountExpirationDate) ||
		(u.CredentialsExpirationDate != other.CredentialsExpirationDate) ||
		(u.UserType != other.UserType) ||
		(u.LastLogin != other.LastLogin) {

		return false
	}

	return true
}

// Validate validates the User struct using the validator package.
// It checks if the struct fields meet the specified validation rules.
func (u *User) Validate() error {
	var v *validator.Validate = validation.GetValidator()

	if err := v.Struct(u); err != nil {
		return err
	}
	return nil
}
