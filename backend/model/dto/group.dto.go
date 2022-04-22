package dto

import "github.com/google/uuid"

type GroupDto struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
