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
	// Retrieve a single role
	Find(id int) (model.Role, error)
	// Create a permission
	CreatePermission(permission *model.Permission) error
	// Delete a permission
	DeletePermission(id int) error
	// Create a role
	// CreateRole(role *model.Role) error
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

// Retrieve a single role
func (r *DBRoleStore) Find(id int) (model.Role, error) {
	var role model.Role
	err := r.db.First(&role, id).Error
	return role, err
}

// Create permission
//
// FIXME: should this be in its own store?
func (r *DBRoleStore) CreatePermission(permission *model.Permission) error {
	err := r.db.Create(permission).Error
	return err
}

// Delete a permission
func (r *DBRoleStore) DeletePermission(id int) error {
	return r.db.Delete(&model.Permission{}, id).Error
}
