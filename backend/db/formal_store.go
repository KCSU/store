package db

import (
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

// TODO: return HTTP errors instead

type FormalStore struct {
	db *gorm.DB
}

func NewFormalStore(db *gorm.DB) *FormalStore {
	return &FormalStore{
		db: db,
	}
}

func (f *FormalStore) Get() ([]model.Formal, error) {
	var data []model.Formal
	err := f.db.Where("date_time > NOW()").Find(&data).Error
	return data, err
}

func (f *FormalStore) GetWithQueue() ([]model.Formal, error) {
	var data []model.Formal
	err := f.db.Preload("TicketSales", "is_queue").Where("date_time > NOW()").Find(&data).Error
	return data, err
}

func (f *FormalStore) All() ([]model.Formal, error) {
	var data []model.Formal
	err := f.db.Find(&data).Error
	return data, err
}

func (f *FormalStore) Find(id int) (model.Formal, error) {
	var formal model.Formal
	err := f.db.First(&formal, id).Error
	return formal, err
}

func (f *FormalStore) TicketsRemaining(formal *model.Formal, isGuest bool) uint {
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
