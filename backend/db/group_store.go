package db

import (
	"github.com/google/uuid"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Helper for using Groups in the database
type GroupStore interface {
	// Retrieve all groups
	Get() ([]model.Group, error)
	// Retrieve a single group
	Find(id uuid.UUID) (model.Group, error)
	// Add a user to the group
	AddUser(group *model.Group, email string) error
	// Remove a user from the group
	RemoveUser(group *model.Group, email string) error
	// Replace the lookup users in the group
	ReplaceLookupUsers(group *model.Group, users []model.GroupUser) error
	// Create a group
	Create(group *model.Group) error
	// Update a group
	Update(group *model.Group) error
	// Delete a group
	Delete(group *model.Group) error
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
func (g *DBGroupStore) Find(id uuid.UUID) (model.Group, error) {
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

// Replace the lookup users in the group
func (g *DBGroupStore) ReplaceLookupUsers(group *model.Group, users []model.GroupUser) error {
	// err := g.db.Model(&group).
	// 	Not("is_manual").Association("GroupUsers").Replace(&users)
	newUsers := make([]model.GroupUser, len(users))
	for i, user := range users {
		newUsers[i] = model.GroupUser{
			GroupID:   group.ID,
			UserEmail: user.UserEmail,
			IsManual:  false,
		}
	}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("group_id", group.ID).Not("is_manual").Delete(&model.GroupUser{}).Error
		if err != nil {
			return err
		}
		if err := tx.Create(&newUsers).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// Create a group
func (g *DBGroupStore) Create(group *model.Group) error {
	return g.db.Create(group).Error
}

// Update a group
func (g *DBGroupStore) Update(group *model.Group) error {
	return g.db.Omit("created_at").Save(group).Error
}

// Delete a group
func (g *DBGroupStore) Delete(group *model.Group) error {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("group_id = ?", group.ID).
			Delete(&model.GroupUser{}).Error
		if err != nil {
			return err
		}
		if err := tx.Delete(group).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
