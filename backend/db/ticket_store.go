package db

import (
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

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

func (t *TicketStore) Create(tickets []model.Ticket) error {
	return t.db.Create(tickets).Error
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
