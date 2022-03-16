package db

import (
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

// Helper for using Groups in the database
type GroupStore interface {
	// Retrieve all groups
	Get() ([]model.Group, error)
}

// Helper struct for using Groups in the database
type DBGroupStore struct {
	db *gorm.DB
}

func NewGroupStore(db *gorm.DB) GroupStore {
	return &DBGroupStore{
		db: db,
	}
}

// Retrieve all groups
func (g *DBGroupStore) Get() ([]model.Group, error) {
	var data []model.Group
	err := g.db.Find(&data).Error
	return data, err
}
