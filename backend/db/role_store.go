package db

import (
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

// Helper for using Roles in the database
type RoleStore interface {
	// Retrieve all roles with permissions
	Get() ([]model.Role, error)
	// Retrieve user-role mapping
	GetUserRoles() ([]model.UserRole, error)
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

// Retrieve user-role mapping
func (r *DBRoleStore) GetUserRoles() ([]model.UserRole, error) {
	var userRoles []model.UserRole
	err := r.db.Table("user_roles").Joins("User").Joins("Role").Find(&userRoles).Error
	return userRoles, err
}
