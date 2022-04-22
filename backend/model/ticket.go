package model

import "github.com/google/uuid"

type Ticket struct {
	Model
	IsGuest    bool      `json:"isGuest"`
	IsQueue    bool      `json:"isQueue"`
	MealOption string    `json:"option"`
	FormalID   uuid.UUID `json:"formalId"`
	Formal     *Formal   `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	UserID     uuid.UUID `json:"userId"`
	User       *User     `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}
