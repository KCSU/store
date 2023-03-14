package dto

import (
	"time"

	"github.com/google/uuid"
)

type ScannedTicketDto struct {
	ID         uuid.UUID `json:"id"`
	IsGuest    bool      `json:"isGuest"`
	IsScanned  bool      `json:"isScanned"`
	MealOption string    `json:"option"`
	FormalID   uuid.UUID `json:"formalId"`
	FormalName string    `json:"formalName"`
	FormalDate time.Time `json:"formalDate"`
	UserName   string    `json:"userName"`
}
