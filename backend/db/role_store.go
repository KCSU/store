package db

import (
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

// Helper for using Roles in the database
type RoleStore interface {
	// Retrieve all roles with permissions
	Get() ([]model.Role, error)
}

// Helper struct for using Roles in the database
type DBRoleStore struct {
	db *gorm.DB
}

// Initialise the role helper
func NewRoleStore(db *gorm.DB) RoleStore {
	return &DBRoleStore{
		db: db,
	}
}

// Retrieve all roles with permissions
func (r *DBRoleStore) Get() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}
