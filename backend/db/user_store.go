package db

import (
	"github.com/kcsu/store/model"
	"github.com/markbates/goth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

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

func (u *UserStore) Exists(email string) (bool, error) {
	var count int64
	err := u.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
