package dto

type BuyTicketDto struct {
	FormalId     uint            `json:"formalId"`
	Ticket       TicketRequest   `json:"ticket"`
	GuestTickets []TicketRequest `json:"guestTickets"`
}
