package service

import (
	"fmt"

	"github.com/yoanesber/go-consumer-api-with-jwt/config/database"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/repository"
)

// Interface for role service
// This interface defines the methods that the role service should implement
type RoleService interface {
	GetRoleByID(id uint) (entity.Role, error)
	GetRoleByName(name string) (entity.Role, error)
}

// This struct defines the RoleService that contains a repository field of type RoleRepository
// It implements the RoleService interface and provides methods for role-related operations
type roleService struct {
	repo repository.RoleRepository
}

// NewRoleService creates a new instance of RoleService with the given repository.
// It initializes the roleService struct and returns it.
func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

// GetRoleByID retrieves a role by its ID from the database.
func (s *roleService) GetRoleByID(id uint) (entity.Role, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.Role{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the role by ID from the repository
	role, err := s.repo.GetRoleByID(db, id)
	if err != nil {
		return entity.Role{}, err
	}

	return role, nil
}

// GetRoleByName retrieves a role by its name from the database.
func (s *roleService) GetRoleByName(name string) (entity.Role, error) {
	db := database.GetPostgres()
	if db == nil {
		return entity.Role{}, fmt.Errorf("database connection is nil")
	}

	// Retrieve the role by name from the repository
	role, err := s.repo.GetRoleByName(db, name)
	if err != nil {
		return entity.Role{}, err
	}

	return role, nil
}
