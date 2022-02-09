package dto

type TicketRequestDto struct {
	MealOption string `json:"option" validate:"required,min=3,max=100"`
}
