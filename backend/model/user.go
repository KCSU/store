package model

// TODO: not nulls?
type User struct {
	Model
	Name           string `json:"name"`
	Email          string `json:"email" gorm:"uniqueIndex;not null"`
	ProviderUserId string `json:"-" gorm:"uniqueIndex; not null"`
	AccessToken    string `json:"-"`
	RefreshToken   string `json:"-"`
}
