package model

type Role struct {
	Model
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions,omitempty"`
}

// FIXME: should this be registered as an official join table?
type UserRole struct {
	UserID uint
	User   User
	RoleID uint
	Role   Role
}
