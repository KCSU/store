package dto

type AddExtraToBillDto struct {
	Description string  `json:"description" validate:"required"`
	Amount      float32 `json:"amount" validate:"required"`
}
