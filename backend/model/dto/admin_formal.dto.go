package dto

import "github.com/kcsu/store/model"

type AdminFormalDto struct {
	model.Formal
	Groups        []GroupDto           `json:"groups"`
	TicketSales   []AdminTicketDto     `json:"ticketSales"`
	ManualTickets []model.ManualTicket `json:"manualTickets"`
}
