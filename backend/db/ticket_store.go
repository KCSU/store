package db

import (
	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"gorm.io/gorm"
)

// TODO: return HTTP errors instead

type TicketStore struct {
	db *gorm.DB
}

func NewTicketStore(db *gorm.DB) *TicketStore {
	return &TicketStore{
		db: db,
	}
}

func (t *TicketStore) Get() ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := t.db.Preload("Formal").Find(&tickets).Error
	return tickets, err
}

func (t *TicketStore) Find(id int) (model.Ticket, error) {
	var ticket model.Ticket
	err := t.db.First(&ticket, id).Error
	return ticket, err
}

// TODO: rewrite completely
func (t *TicketStore) Update(id int, ticket *dto.TicketRequestDto) error {
	return t.db.Model(&model.Ticket{}).Where("id = ?", id).Update("meal_option", ticket.MealOption).Error
}

func (t *TicketStore) BatchCreate(tickets []model.Ticket) error {
	return t.db.Create(tickets).Error
}

func (t *TicketStore) Create(ticket *model.Ticket) error {
	return t.db.Create(ticket).Error
}

func (t *TicketStore) CountGuestByFormal(formalID int) (int64, error) {
	var count int64
	err := t.db.Model(&model.Ticket{}).Where("is_guest").Where("formal_id = ?", formalID).Count(&count).Error
	return count, err
}

func (t *TicketStore) ExistsByFormal(formalID int) (bool, error) {
	var count int64
	err := t.db.Model(&model.Ticket{}).Not("is_guest").Where("formal_id = ?", formalID).Count(&count).Error
	return count > 0, err
}

func (t *TicketStore) DeleteByFormal(formalID int) error {
	return t.db.Where("formal_id = ?", formalID).Delete(&model.Ticket{}).Error
}

func (t *TicketStore) Delete(id int) error {
	return t.db.Delete(&model.Ticket{}, id).Error
}
