package model

import (
	"time"

	"github.com/google/uuid"
)

type FormalCostBreakdown struct {
	FormalID       uuid.UUID `json:"formalId"`
	Name           string    `json:"formalName"`
	Price          float32   `json:"price"`
	GuestPrice     float32   `json:"guestPrice"`
	DateTime       time.Time `json:"dateTime"`
	Standard       int       `json:"standard"`
	Guest          int       `json:"guest"`
	StandardManual int       `json:"standardManual"`
	GuestManual    int       `json:"guestManual"`
}

type UserCostBreakdown struct {
	Email string  `json:"userEmail"`
	Cost  float32 `json:"cost"`
}
