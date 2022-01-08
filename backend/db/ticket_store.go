package db

import (
	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"gorm.io/gorm"
)

// TODO: return HTTP errors instead
// TODO: DATE CHECKS

// Helper struct for using Tickets in the database
type TicketStore struct {
	db *gorm.DB
}

// Initialise the ticket helper
func NewTicketStore(db *gorm.DB) *TicketStore {
	return &TicketStore{
		db: db,
	}
}

// Retrieve all the user's tickets
func (t *TicketStore) Get(userId int) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := t.db.Preload("Formal").Where("user_id = ?", userId).Find(&tickets).Error
	return tickets, err
}

// Get a ticket by its id
func (t *TicketStore) Find(id int) (model.Ticket, error) {
	var ticket model.Ticket
	err := t.db.First(&ticket, id).Error
	return ticket, err
}

// Update a ticket meal option by id
//
// TODO: rewrite completely
func (t *TicketStore) Update(id int, ticket *dto.TicketRequestDto) error {
	return t.db.Model(&model.Ticket{}).Where("id = ?", id).Update("meal_option", ticket.MealOption).Error
}

// Create tickets in batch
func (t *TicketStore) BatchCreate(tickets []model.Ticket) error {
	return t.db.Create(tickets).Error
}

// Create a single ticket
func (t *TicketStore) Create(ticket *model.Ticket) error {
	return t.db.Create(ticket).Error
}

// Get tne number of guests attending a specified formal
func (t *TicketStore) CountGuestByFormal(formalID int, userID int) (int64, error) {
	var count int64
	err := t.db.Model(&model.Ticket{}).Where("is_guest").Where("formal_id = ? AND user_id = ?", formalID, userID).Count(&count).Error
	return count, err
}

// Check if a user has purchased a ticket for a specified formal
func (t *TicketStore) ExistsByFormal(formalID int, userID int) (bool, error) {
	var count int64
	err := t.db.Model(&model.Ticket{}).Not("is_guest").Where("formal_id = ? AND user_id = ?", formalID, userID).Count(&count).Error
	return count > 0, err
}

// Delete all the user's tickets for a specified formal
func (t *TicketStore) DeleteByFormal(formalID int, userID int) error {
	return t.db.Where("formal_id = ? AND user_id = ?", formalID, userID).Delete(&model.Ticket{}).Error
}

// Delete a single ticket
func (t *TicketStore) Delete(id int) error {
	return t.db.Delete(&model.Ticket{}, id).Error
}
