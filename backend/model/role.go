package model

import "github.com/google/uuid"

type Role struct {
	Model
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions,omitempty"`
	Users       []User       `json:"-" gorm:"many2many:user_roles;"`
}

// HACK: should this be registered as an official join table?
type UserRole struct {
	UserID uuid.UUID
	User   User
	RoleID uuid.UUID
	Role   Role
}
