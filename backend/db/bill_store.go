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
	// Create a bill
	Create(bill *model.Bill) error
	// Update a bill
	Update(bill *model.Bill) error
	// Delete a bill
	Delete(bill *model.Bill) error
	// Add a formal to a bill
	AddFormals(bill *model.Bill, formalIds []uuid.UUID) error
	// Remove a formal from a bill
	RemoveFormal(bill *model.Bill, formalId uuid.UUID) error
	// Get bill cost breakdown by formal
	GetCostBreakdown(bill *model.Bill) ([]model.FormalCostBreakdown, error)
	// Get bill cost breakdown by user
	GetCostBreakdownByUser(bill *model.Bill) ([]model.UserCostBreakdown, error)
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
	err := b.db.Order("start DESC").Find(&data).Error
	return data, err
}

// Create a bill
func (b *DBBillStore) Create(bill *model.Bill) error {
	return b.db.Create(bill).Error
}

// Update a bill
func (b *DBBillStore) Update(bill *model.Bill) error {
	return b.db.Omit("Formals", "created_at").Save(bill).Error
}

// Delete a bill
func (b *DBBillStore) Delete(bill *model.Bill) error {
	return b.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(bill).Association("Formals").Clear()
		if err != nil {
			return err
		}
		return tx.Delete(bill).Error
	})
}

// Add a formal to a bill
func (b *DBBillStore) AddFormals(bill *model.Bill, formalIds []uuid.UUID) error {
	return b.db.Model(&model.Formal{}).Where("id in ?", formalIds).Update("bill_id", bill.ID).Error
}

// Remove a formal from a bill
func (b *DBBillStore) RemoveFormal(bill *model.Bill, formalId uuid.UUID) error {
	f := &model.Formal{}
	f.ID = formalId
	return b.db.Model(f).Update("bill_id", nil).Error
}

func (b *DBBillStore) billQuery(bill *model.Bill) *gorm.DB {
	// Select email, formal ID, guest status from tickets
	tickets := b.db.Model(&model.Ticket{}).
		Not("is_queue").
		Joins("LEFT JOIN users ON users.id = tickets.user_id").
		Select("users.email, tickets.formal_id, tickets.is_guest")
	// Select email, formal ID, guest status from manual tickets
	manualTickets := b.db.Model(&model.ManualTicket{}).
		Where("type IN ('standard', 'ents', 'guest')").
		Select(`(CASE WHEN type = 'ents' THEN 'ents' ELSE manual_tickets.email END) AS email,
		manual_tickets.formal_id AS formal_id,
		manual_tickets.type = 'guest' AS is_guest`)
	// Combine these with a union statement
	allTickets := b.db.Raw(`(?) UNION ALL (?)`, tickets, manualTickets)
	// Query by bill ID and formals
	return b.db.Model(&model.Formal{}).
		Where("bill_id = ?", bill.ID).
		Joins("RIGHT JOIN (?) ticket ON formal_id = formals.id", allTickets)
}

// Get bill cost breakdown by formal
func (b *DBBillStore) GetCostBreakdown(bill *model.Bill) ([]model.FormalCostBreakdown, error) {
	var data []model.FormalCostBreakdown
	err := b.billQuery(bill).Group("formals.id").
		Select(`formals.id as formal_id, formals.name, formals.price, formals.guest_price, formals.date_time,
		SUM(CASE WHEN NOT ticket.is_guest THEN 1 ELSE 0 END) AS standard,
		SUM(CASE WHEN ticket.is_guest THEN 1 ELSE 0 END) AS guest`).
		Order("formals.date_time").
		Scan(&data).Error
	return data, err
}

// Get bill cost breakdown by user
func (b *DBBillStore) GetCostBreakdownByUser(bill *model.Bill) ([]model.UserCostBreakdown, error) {
	var data []model.UserCostBreakdown
	err := b.billQuery(bill).Group("ticket.email").
		Select(`ticket.email,
		SUM(CASE WHEN ticket.is_guest
			THEN formals.guest_price
			ELSE formals.price
		END) AS cost`).
		Order("ticket.email").
		Scan(&data).Error
	return data, err
}
