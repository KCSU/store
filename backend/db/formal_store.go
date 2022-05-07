package db

import (
	"github.com/google/uuid"
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

// TODO: return HTTP errors instead

// Helper for using Formals in the database
type FormalStore interface {
	// Retrieve all upcoming formals
	Get() ([]model.Formal, error)
	// Retrieve all upcoming formals with groups
	GetWithUserData(userId uuid.UUID) ([]model.Formal, error)
	// Retrieve all formals
	All() ([]model.Formal, error)
	// Get a formal by id
	Find(id uuid.UUID) (model.Formal, error)
	// Get a formal by id with tickets
	FindWithTickets(id uuid.UUID) (model.Formal, error)
	// Get the guest list for a formal
	FindGuestList(id uuid.UUID) ([]model.FormalGuest, error)
	// Get the number of tickets remaining for a specified formal
	TicketsRemaining(formal *model.Formal, isGuest bool) uint
	// Create a formal
	Create(formal *model.Formal) error
	// Find all groups with specified ids
	GetGroups(ids []uuid.UUID) ([]model.Group, error)
	// Update a formal
	Update(formal *model.Formal) error
	// Delete a formal
	Delete(formal *model.Formal) error
	// Update groups for a formal
	UpdateGroups(formal model.Formal, groups []model.Group) error
}

// Helper struct for using Formals in the database
type DBFormalStore struct {
	db *gorm.DB
}

// Initialise the formal helper
func NewFormalStore(db *gorm.DB) FormalStore {
	return &DBFormalStore{
		db: db,
	}
}

// Retrieve all upcoming formals
func (f *DBFormalStore) Get() ([]model.Formal, error) {
	var data []model.Formal
	err := f.db.Where("date_time > NOW()").Order("date_time").Find(&data).Error
	return data, err
}

// Retrieve all upcoming formals with groups
func (f *DBFormalStore) GetWithUserData(userId uuid.UUID) ([]model.Formal, error) {
	var data []model.Formal
	err := f.db.Where("date_time > NOW()").Order("date_time").
		Preload("Groups").
		Preload("TicketSales", "user_id = ?", userId).
		Find(&data).Error
	return data, err
}

// Retrieve all formals
func (f *DBFormalStore) All() ([]model.Formal, error) {
	var data []model.Formal
	err := f.db.Order("date_time DESC").Find(&data).Error
	return data, err
}

// Get a formal by id
func (f *DBFormalStore) Find(id uuid.UUID) (model.Formal, error) {
	var formal model.Formal
	err := f.db.Preload("Groups").First(&formal, id).Error
	return formal, err
}

// Get a formal by id with tickets
func (f *DBFormalStore) FindWithTickets(id uuid.UUID) (model.Formal, error) {
	var formal model.Formal
	err := f.db.Preload("Groups").
		Preload("ManualTickets").
		Preload("TicketSales", "NOT is_queue").
		Preload("TicketSales.User").First(&formal, id).Error
	return formal, err
}

// Get the guest list for a formal
func (f *DBFormalStore) FindGuestList(id uuid.UUID) ([]model.FormalGuest, error) {
	var guests []model.FormalGuest
	err := f.db.Model(&model.Ticket{}).
		Where("formal_id = ?", id).
		Not("is_queue").
		Joins("LEFT JOIN users ON users.id = tickets.user_id").
		Group("users.id").
		Order("users.name").
		Select(`users.name, users.email,
			COUNT(*) FILTER (WHERE tickets.is_guest) AS guests`).
		Scan(&guests).Error
	return guests, err
}

// Get the number of tickets remaining for a specified formal
func (f *DBFormalStore) TicketsRemaining(formal *model.Formal, isGuest bool) uint {
	var query string
	var baseTickets uint
	if isGuest {
		baseTickets = formal.GuestTickets
		query = "is_guest AND NOT is_queue"
	} else {
		baseTickets = formal.Tickets
		query = "NOT is_guest AND NOT is_queue"
	}

	return baseTickets - uint(
		f.db.Model(formal).Where(query).Association("TicketSales").Count(),
	)
}

// Find all groups with specified ids
func (f *DBFormalStore) GetGroups(ids []uuid.UUID) ([]model.Group, error) {
	if len(ids) == 0 {
		return []model.Group{}, nil
	}
	var groups []model.Group
	err := f.db.Find(&groups, ids).Error
	return groups, err
}

// Create a formal
func (f *DBFormalStore) Create(formal *model.Formal) error {
	err := f.db.Omit("Groups.*").Create(formal).Error
	return err
}

// Update a formal
func (f *DBFormalStore) Update(formal *model.Formal) error {
	return f.db.Omit("CreatedAt").Save(formal).Error
}

// Update groups for a formal
func (f *DBFormalStore) UpdateGroups(formal model.Formal, groups []model.Group) error {
	return f.db.Model(&formal).
		Omit("Groups.*").Association("Groups").
		Replace(&groups)
}

// Delete a formal
func (f *DBFormalStore) Delete(formal *model.Formal) error {
	return f.db.Delete(formal).Error
}
