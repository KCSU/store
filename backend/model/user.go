package model

// TODO: not nulls?
type User struct {
	Model
	Name           string `json:"name"`
	Email          string `json:"email" gorm:"unique;uniqueIndex;not null"`
	ProviderUserId string `json:"-" gorm:"unique;uniqueIndex; not null"`
	Roles          []Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}
