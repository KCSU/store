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
	err := b.db.Find(&data).Error
	return data, err
}

// Update a bill
func (b *DBBillStore) Update(bill *model.Bill) error {
	return b.db.Omit("Formals", "created_at").Save(bill).Error
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

// Get bill cost breakdown by formal
func (b *DBBillStore) GetCostBreakdown(bill *model.Bill) ([]model.FormalCostBreakdown, error) {
	var data []model.FormalCostBreakdown
	// This is a really complicated query. Let me explain it:
	// We want to get the breakdown of the bill by formal.
	// So we get the list of formals associated with the bill.
	// Then, we get the list of tickets & manual tickets associated with those formals;
	// We sum up the cost of each ticket, grouped by formal ID and ticket type.
	err := b.db.Raw(`
		SELECT
			formal.id as formal_id, formal.name, formal.price, formal.guest_price, formal.date_time,
			SUM(CASE WHEN NOT ticket.is_guest THEN 1 ELSE 0 END) AS standard,
			SUM(CASE WHEN ticket.is_guest THEN 1 ELSE 0 END) AS guest,
			SUM(CASE WHEN manual_ticket.type IN ('standard', 'ents') THEN 1 ELSE 0 END) AS standard_manual,
			SUM(CASE WHEN manual_ticket.type = 'guest' THEN 1 ELSE 0 END) AS guest_manual
		FROM
			formals AS formal
			LEFT JOIN tickets AS ticket ON formal.id = ticket.formal_id
				AND NOT ticket.is_queue
				AND ticket.deleted_at IS NULL
			LEFT JOIN manual_tickets AS manual_ticket ON formal.id = manual_ticket.formal_id
				AND manual_ticket.deleted_at IS NULL
		WHERE
			formal.bill_id = ?
		GROUP BY
			formal.id
	`, bill.ID).Scan(&data).Error
	return data, err
}

// Get bill cost breakdown by user
func (b *DBBillStore) GetCostBreakdownByUser(bill *model.Bill) ([]model.UserCostBreakdown, error) {
	var data []model.UserCostBreakdown
	// This is a really complicated query. Let me explain it:
	// We want to get the breakdown of the bill by user.
	// So we get the list of formals associated with the bill.
	// Then, we get the list of tickets & manual tickets associated with those formals;
	// using the UNION ALL to combine the two lists.
	// Finally, we sum up the cost of each ticket, grouped by user email, using is_guest
	// to assess whether the ticket is a guest or not.
	err := b.db.Raw(`
		SELECT
			ticket.email,
			SUM(CASE WHEN ticket.is_guest THEN formal.guest_price ELSE formal.price END) AS cost
		FROM
			formals AS formal
			RIGHT JOIN (
				SELECT
					users.email AS email,
					tickets.formal_id AS formal_id,
					tickets.is_guest AS is_guest
				FROM
					tickets
					LEFT JOIN users ON users.id = tickets.user_id
				WHERE
					NOT is_queue AND tickets.deleted_at IS NULL
				UNION ALL
				SELECT
					(CASE WHEN type = 'ents' THEN 'ents' ELSE manual_tickets.email END) AS email,
					manual_tickets.formal_id AS formal_id,
					manual_tickets.type = 'guest' AS is_guest
				FROM manual_tickets
				WHERE
					deleted_at IS NULL
					AND type IN ('standard', 'ents', 'guest')
			) ticket ON formal.id = ticket.formal_id
		WHERE
			formal.bill_id = ?
		GROUP BY
			ticket.email
	`, bill.ID).Scan(&data).Error
	return data, err
}
