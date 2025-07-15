package repository

import (
	"gorm.io/gorm"

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
)

// Interface for role repository
// This interface defines the methods that the role repository should implement
type RoleRepository interface {
	GetRoleByID(tx *gorm.DB, id uint) (entity.Role, error)
	GetRoleByName(tx *gorm.DB, name string) (entity.Role, error)
}

// This struct defines the RoleRepository that contains methods for interacting with the database
type roleRepository struct{}

// NewRoleRepository creates a new instance of RoleRepository.
// It initializes the roleRepository struct and returns it.
func NewRoleRepository() RoleRepository {
	return &roleRepository{}
}

// GetRoleByID retrieves a role by its ID from the database.
func (r *roleRepository) GetRoleByID(tx *gorm.DB, id uint) (entity.Role, error) {
	// Select the role with the given ID from the database
	var role entity.Role
	err := tx.First(&role, "id = ?", id).Error

	if err != nil {
		return entity.Role{}, err
	}

	return role, nil
}

// GetRoleByName retrieves a role by its name from the database.
func (r *roleRepository) GetRoleByName(tx *gorm.DB, name string) (entity.Role, error) {
	// Select the role with the given name from the database
	var role entity.Role
	err := tx.First(&role, "lower(name) = lower(?)", name).Error

	if err != nil {
		return entity.Role{}, err
	}

	return role, nil
}
