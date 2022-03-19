package db

import (
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Helper for using Groups in the database
type GroupStore interface {
	// Retrieve all groups
	Get() ([]model.Group, error)
	// Retrieve a single group
	Find(id int) (model.Group, error)
	// Add a user to the group
	AddUser(group *model.Group, email string) error
	// Remove a user from the group
	RemoveUser(group *model.Group, email string) error
	// Create a group
	Create(group *model.Group) error
	// Update a group
	Update(group *model.Group) error
}

// Helper struct for using Groups in the database
type DBGroupStore struct {
	db *gorm.DB
}

func NewGroupStore(db *gorm.DB) GroupStore {
	return &DBGroupStore{
		db: db,
	}
}

// Retrieve all groups
func (g *DBGroupStore) Get() ([]model.Group, error) {
	var data []model.Group
	err := g.db.Find(&data).Error
	return data, err
}

// Retrieve a single group
func (g *DBGroupStore) Find(id int) (model.Group, error) {
	var group model.Group
	err := g.db.Preload("GroupUsers").First(&group, id).Error
	return group, err
}

// Add a user to the group
func (g *DBGroupStore) AddUser(group *model.Group, email string) error {
	groupUser := model.GroupUser{
		UserEmail: email,
		IsManual:  true,
	}
	err := g.db.Model(group).Association("GroupUsers").Append(&groupUser)
	return err
}

// Remove a user from the group
func (g *DBGroupStore) RemoveUser(group *model.Group, email string) error {
	res := g.db.Where("group_id = ? AND user_email = ?", group.ID, email).
		Where("is_manual").Delete(&model.GroupUser{})
	if res.RowsAffected == 0 {
		return echo.ErrNotFound
	}
	return res.Error
}

// Create a group
func (g *DBGroupStore) Create(group *model.Group) error {
	return g.db.Create(group).Error
}

// Update a group
func (g *DBGroupStore) Update(group *model.Group) error {
	return g.db.Omit("created_at").Save(group).Error
}
