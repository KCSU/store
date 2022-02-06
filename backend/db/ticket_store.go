package db

import (
	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"gorm.io/gorm"
)

// TODO: return HTTP errors instead
// TODO: DATE CHECKS

// Helper for using Tickets in the database
type TicketStore interface {
	// Retrieve all the user's tickets
	Get(userId int) ([]model.Ticket, error)
	// Get a ticket by its id
	Find(id int) (model.Ticket, error)
	// Get a ticket by its id with formal relation
	FindWithFormal(id int) (model.Ticket, error)
	// Update a ticket meal option by id
	Update(id int, ticket *dto.TicketRequestDto) error
	// Create tickets in batch
	BatchCreate(tickets []model.Ticket) error
	// Create a single ticket
	Create(ticket *model.Ticket) error
	// Get tne number of guests attending a specified formal
	CountGuestByFormal(formalId int, userId int) (int64, error)
	// Delete all the user's tickets for a specified formal
	ExistsByFormal(formalID int, userID int) (bool, error)
	// Delete all the user's tickets for a specified formal
	DeleteByFormal(formalID int, userID int) error
	// Delete a single ticket
	Delete(id int) error
}

// Helper struct for using Tickets in the database
type DBTicketStore struct {
	db *gorm.DB
}

// Initialise the ticket helper
func NewTicketStore(db *gorm.DB) TicketStore {
	return &DBTicketStore{
		db: db,
	}
}

// Retrieve all the user's tickets
func (t *DBTicketStore) Get(userId int) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := t.db.Preload("Formal").Where("user_id = ?", userId).Find(&tickets).Error
	return tickets, err
}

// Get a ticket by its id
func (t *DBTicketStore) Find(id int) (model.Ticket, error) {
	var ticket model.Ticket
	err := t.db.First(&ticket, id).Error
	return ticket, err
}

// Get a ticket by its id
func (t *DBTicketStore) FindWithFormal(id int) (model.Ticket, error) {
	var ticket model.Ticket
	err := t.db.Preload("Formal").First(&ticket, id).Error
	return ticket, err
}

// Update a ticket meal option by id
//
// TODO: rewrite completely
func (t *DBTicketStore) Update(id int, ticket *dto.TicketRequestDto) error {
	return t.db.Model(&model.Ticket{}).Where("id = ?", id).Update("meal_option", ticket.MealOption).Error
}

// Create tickets in batch
func (t *DBTicketStore) BatchCreate(tickets []model.Ticket) error {
	return t.db.Create(tickets).Error
}

// Create a single ticket
func (t *DBTicketStore) Create(ticket *model.Ticket) error {
	return t.db.Create(ticket).Error
}

// Get tne number of guests attending a specified formal
func (t *DBTicketStore) CountGuestByFormal(formalID int, userID int) (int64, error) {
	var count int64
	err := t.db.Model(&model.Ticket{}).Where("is_guest").Where("formal_id = ? AND user_id = ?", formalID, userID).Count(&count).Error
	return count, err
}

// Check if a user has purchased a ticket for a specified formal
func (t *DBTicketStore) ExistsByFormal(formalID int, userID int) (bool, error) {
	var count int64
	err := t.db.Model(&model.Ticket{}).Not("is_guest").Where("formal_id = ? AND user_id = ?", formalID, userID).Count(&count).Error
	return count > 0, err
}

// Delete all the user's tickets for a specified formal
func (t *DBTicketStore) DeleteByFormal(formalID int, userID int) error {
	return t.db.Where("formal_id = ? AND user_id = ?", formalID, userID).Delete(&model.Ticket{}).Error
}

// Delete a single ticket
func (t *DBTicketStore) Delete(id int) error {
	return t.db.Delete(&model.Ticket{}, id).Error
}
