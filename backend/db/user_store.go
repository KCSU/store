package db

import (
	"github.com/kcsu/store/model"
	"github.com/markbates/goth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Helper struct for Users in the database
type UserStore struct {
	db *gorm.DB
}

// Initialise the user helper
func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// Update, retrieve or create a user from OAuth data
func (u *UserStore) FindOrCreate(gu goth.User) (model.User, error) {
	user := model.User{
		Name:           gu.Name,
		Email:          gu.Email,
		ProviderUserId: gu.UserID,
		AccessToken:    gu.AccessToken,
		RefreshToken:   gu.RefreshToken,
	}
	err := u.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider_user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"access_token", "refresh_token"}),
	}).Create(&user).Error
	return user, err
}

func (u *UserStore) Find(id int) (model.User, error) {
	var user model.User
	err := u.db.First(&user, id).Error
	return user, err
}

// Check if a user with the specified email exists
func (u *UserStore) Exists(email string) (bool, error) {
	var count int64
	err := u.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (u *UserStore) Groups(user *model.User) ([]model.Group, error) {
	var groups []model.Group
	err := u.db.Joins("inner join group_users on id = group_users.group_id").Find(&groups, "group_users.user_email = ?", user.Email).Error
	return groups, err
}
