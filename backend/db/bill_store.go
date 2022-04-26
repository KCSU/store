package db

import (
	"github.com/google/uuid"
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

// Helper for using Bills in the database
type BillStore interface {
	// Retrieve all bills
	Get() ([]model.Bill, error)
	// Retrieve a single bill
	Find(id uuid.UUID) (model.Bill, error)
	// Retrieve a single bill with formals
	FindWithFormals(id uuid.UUID) (model.Bill, error)
	// Update a bill
	Update(bill *model.Bill) error
	// Add a formal to a bill
	AddFormal(bill *model.Bill, formalId uuid.UUID) error
	// Remove a formal from a bill
	RemoveFormal(bill *model.Bill, formalId uuid.UUID) error
}

// Helper struct for using Bills in the database
type DBBillStore struct {
	db *gorm.DB
}

func NewBillStore(db *gorm.DB) BillStore {
	return &DBBillStore{
		db: db,
	}
}

// Retrieve a single bill
func (b *DBBillStore) Find(id uuid.UUID) (model.Bill, error) {
	var bill model.Bill
	err := b.db.First(&bill, id).Error
	return bill, err
}

// Retrieve a single bill with formals
func (b *DBBillStore) FindWithFormals(id uuid.UUID) (model.Bill, error) {
	var bill model.Bill
	err := b.db.Preload("Formals").First(&bill, id).Error
	return bill, err
}

// Retrieve all bills
func (b *DBBillStore) Get() ([]model.Bill, error) {
	var data []model.Bill
	err := b.db.Find(&data).Error
	return data, err
}

// Update a bill
func (b *DBBillStore) Update(bill *model.Bill) error {
	return b.db.Omit("Formals", "created_at").Save(bill).Error
}

// Add a formal to a bill
func (b *DBBillStore) AddFormal(bill *model.Bill, formalId uuid.UUID) error {
	f := &model.Formal{}
	f.ID = formalId
	return b.db.Model(f).Update("bill_id", bill.ID).Error
}

// Remove a formal from a bill
func (b *DBBillStore) RemoveFormal(bill *model.Bill, formalId uuid.UUID) error {
	f := &model.Formal{}
	f.ID = formalId
	return b.db.Model(f).Update("bill_id", nil).Error
}
