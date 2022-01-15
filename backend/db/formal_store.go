package db

import (
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

// TODO: return HTTP errors instead

// Helper for using Formals in the database
type FormalStore interface {
	// Retrieve all upcoming formals
	Get() ([]model.Formal, error)
	// Retrieve all formals
	All() ([]model.Formal, error)
	// Get a formal by id
	Find(id int) (model.Formal, error)
	// Get the number of tickets remaining for a specified formal
	TicketsRemaining(formal *model.Formal, isGuest bool) uint
}

// Helper struct for using Formals in the database
type DBFormalStore struct {
	db *gorm.DB
}

// Initialise the formal helper
func NewFormalStore(db *gorm.DB) FormalStore {
	return &DBFormalStore{
		db: db,
	}
}

// Retrieve all upcoming formals
func (f *DBFormalStore) Get() ([]model.Formal, error) {
	var data []model.Formal
	err := f.db.Where("date_time > NOW()").Find(&data).Error
	return data, err
}

// Retrieve all formals
func (f *DBFormalStore) All() ([]model.Formal, error) {
	var data []model.Formal
	err := f.db.Find(&data).Error
	return data, err
}

// Get a formal by id
func (f *DBFormalStore) Find(id int) (model.Formal, error) {
	var formal model.Formal
	err := f.db.Preload("Groups").First(&formal, id).Error
	return formal, err
}

// Get the number of tickets remaining for a specified formal
func (f *DBFormalStore) TicketsRemaining(formal *model.Formal, isGuest bool) uint {
	var query string
	if isGuest {
		query = "is_guest AND NOT is_queue"
	} else {
		query = "NOT is_guest AND NOT is_queue"
	}
	return formal.Tickets - uint(
		f.db.Model(formal).Where(query).Association("TicketSales").Count(),
	)
}
