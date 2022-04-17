package db

import (
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

// Helper for using Manual Tickets in the database
type ManualTicketStore interface {
	// Get a manual ticket by id
	Find(id int) (model.ManualTicket, error)
	// Create a manual ticket
	Create(ticket *model.ManualTicket) error
	// Update a manual ticket
	Update(id int, ticket *model.ManualTicket) error
	// Delete a manual ticket
	Delete(id int) error
}

// Helper struct for using Manual Tickets in the database
type DBManualTicketStore struct {
	db *gorm.DB
}

// Initialise the manual ticket helper
func NewManualTicketStore(db *gorm.DB) ManualTicketStore {
	return &DBManualTicketStore{
		db: db,
	}
}

// Find a manual ticket by id
func (t *DBManualTicketStore) Find(id int) (model.ManualTicket, error) {
	var ticket model.ManualTicket
	err := t.db.First(&ticket, id).Error
	return ticket, err
}

// Create a manual ticket
func (t *DBManualTicketStore) Create(ticket *model.ManualTicket) error {
	return t.db.Create(ticket).Error
}

// Delete a manual ticket
func (t *DBManualTicketStore) Delete(id int) error {
	panic("unimplemented")
}

// Update a manual ticket
func (t *DBManualTicketStore) Update(id int, ticket *model.ManualTicket) error {
	panic("unimplemented")
}
