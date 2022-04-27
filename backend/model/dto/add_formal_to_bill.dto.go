package dto

import "github.com/google/uuid"

type AddFormalToBillDto struct {
	FormalIDs []uuid.UUID `json:"formalIds" validate:"required,dive,required"`
}
