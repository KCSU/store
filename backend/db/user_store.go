package db

import (
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Helper for Users in the database
type UserStore interface {
	// Update, retrieve or create a user from OAuth data
	FindOrCreate(au *auth.OauthUser) (model.User, error)
	// Get a user by id
	Find(id int) (model.User, error)
	// Check if a user with the specified email exists
	Exists(email string) (bool, error)
	// List a user's groups
	Groups(user *model.User) ([]model.Group, error)
}

// Helper struct for Users in the database
type DBUserStore struct {
	db *gorm.DB
}

// Initialise the user helper
func NewUserStore(db *gorm.DB) UserStore {
	return &DBUserStore{
		db: db,
	}
}

// Update, retrieve or create a user from OAuth data
func (u *DBUserStore) FindOrCreate(au *auth.OauthUser) (model.User, error) {
	user := model.User{
		Name:           au.Name,
		Email:          au.Email,
		ProviderUserId: au.UserID,
	}
	err := u.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider_user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&user).Error
	return user, err
}

// Get a user by id
func (u *DBUserStore) Find(id int) (model.User, error) {
	var user model.User
	err := u.db.First(&user, id).Error
	return user, err
}

// Check if a user with the specified email exists
func (u *DBUserStore) Exists(email string) (bool, error) {
	var count int64
	err := u.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// List a user's groups
func (u *DBUserStore) Groups(user *model.User) ([]model.Group, error) {
	var groups []model.Group
	err := u.db.Joins("inner join group_users on id = group_users.group_id").Find(&groups, "group_users.user_email = ?", user.Email).Error
	return groups, err
}
