package model

import "github.com/google/uuid"

type ManualTicket struct {
	Model
	MealOption string    `json:"option"`
	FormalID   uuid.UUID `json:"formalId"`
	Formal     *Formal   `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	// TODO: should this be an enum/custom type?
	// One of: "complimentary", "ents", "standard", "guest"
	Type          string `json:"type"`
	Name          string `json:"name"`
	Justification string `json:"justification"`
	Email         string `json:"email"`
}
