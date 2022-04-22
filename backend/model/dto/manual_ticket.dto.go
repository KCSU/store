package dto

import "github.com/google/uuid"

type CreateManualTicketDto struct {
	MealOption    string    `json:"option" validate:"required,min=3,max=100"`
	FormalID      uuid.UUID `json:"formalId"`
	Type          string    `json:"type" validate:"oneof=complimentary ents standard guest"`
	Name          string    `json:"name" validate:"required,min=3,max=100"`
	Justification string    `json:"justification" validate:"required,min=3,max=100"`
	Email         string    `json:"email" validate:"required,email"`
}

type EditManualTicketDto struct {
	MealOption    string `json:"option" validate:"required,min=3,max=100"`
	Type          string `json:"type" validate:"oneof=complimentary ents standard guest"`
	Name          string `json:"name" validate:"required,min=3,max=100"`
	Justification string `json:"justification" validate:"required,min=3,max=100"`
	Email         string `json:"email" validate:"required,email"`
}
