package dto

type ManualTicketDto struct {
	MealOption    string `json:"option" validate:"required,min=3,max=100"`
	FormalID      int    `json:"formalId"`
	Type          string `json:"type" validate:"oneof=complimentary ents standard guest"`
	Name          string `json:"name" validate:"required,min=3,max=100"`
	Justification string `json:"justification" validate:"required,min=3,max=100"`
	Email         string `json:"email" validate:"required,email"`
}
