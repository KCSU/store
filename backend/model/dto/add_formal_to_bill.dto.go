package dto

import "github.com/google/uuid"

type AddFormalToBillDto struct {
	FormalID uuid.UUID `json:"formalId" validate:"required"`
}
