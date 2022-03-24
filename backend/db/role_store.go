package db

import (
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

type RoleStore interface {
}

type DBRoleStore struct {
	db *gorm.DB
}

func NewRbacStore(db *gorm.DB) RoleStore {
	return &DBRoleStore{
		db: db,
	}
}

func (r *DBRoleStore) Get() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}
