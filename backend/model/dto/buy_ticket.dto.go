package dto

type BuyTicketDto struct {
	FormalId     uint               `json:"formalId"`
	Ticket       TicketRequestDto   `json:"ticket"`
	GuestTickets []TicketRequestDto `json:"guestTickets"`
}
