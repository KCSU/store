package model

import "github.com/google/uuid"

type ExtraCharge struct {
	Model
	Description string    `json:"description"`
	Amount      float32   `json:"amount"`
	BillID      uuid.UUID `json:"billId"`
	Bill        *Bill     `json:"bill,omitempty"`
}
