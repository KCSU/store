package db

import (
	"github.com/google/uuid"
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

// Helper for using Manual Tickets in the database
type ManualTicketStore interface {
	// Get a manual ticket by id
	Find(id uuid.UUID) (model.ManualTicket, error)
	// Create a manual ticket
	Create(ticket *model.ManualTicket) error
	// Update a manual ticket
	Update(ticket *model.ManualTicket) error
	// Delete a manual ticket
	Delete(id uuid.UUID) error
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
func (t *DBManualTicketStore) Find(id uuid.UUID) (model.ManualTicket, error) {
	var ticket model.ManualTicket
	err := t.db.First(&ticket, id).Error
	return ticket, err
}

// Create a manual ticket
func (t *DBManualTicketStore) Create(ticket *model.ManualTicket) error {
	return t.db.Create(ticket).Error
}

// Delete a manual ticket
func (t *DBManualTicketStore) Delete(id uuid.UUID) error {
	return t.db.Delete(&model.ManualTicket{}, id).Error
}

// Update a manual ticket
func (t *DBManualTicketStore) Update(ticket *model.ManualTicket) error {
	return t.db.Omit("CreatedAt").Save(ticket).Error
}
