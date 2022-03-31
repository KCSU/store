package dto

import "github.com/kcsu/store/model"

type AdminTicketDto struct {
	model.Ticket
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
}
