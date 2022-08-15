package model

import (
	"time"

	"github.com/google/uuid"
)

type Formal struct {
	Model
	Name                   string         `json:"name"`
	Menu                   string         `json:"menu"`
	Price                  float32        `json:"price" gorm:"type:decimal(10,2)"`
	GuestPrice             float32        `json:"guestPrice" gorm:"type:decimal(10,2)"`
	GuestLimit             uint           `json:"guestLimit"`
	FirstSaleTickets       uint           `json:"firstSaleTickets"`
	FirstSaleGuestTickets  uint           `json:"firstSaleGuestTickets"`
	FirstSaleStart         time.Time      `json:"firstSaleStart"`
	SecondSaleTickets      uint           `json:"secondSaleTickets"`
	SecondSaleGuestTickets uint           `json:"secondSaleGuestTickets"`
	SecondSaleStart        time.Time      `json:"secondSaleStart"`
	SaleEnd                time.Time      `json:"saleEnd"`
	DateTime               time.Time      `json:"dateTime"`
	HasGuestList           bool           `json:"hasGuestList"`
	IsVisible              bool           `json:"isVisible"`
	BillID                 *uuid.UUID     `json:"billId,omitempty"`
	Bill                   *Bill          `json:"bill,omitempty"`
	TicketSales            []Ticket       `json:"-"`
	ManualTickets          []ManualTicket `json:"-"`
	Groups                 []Group        `json:"groups,omitempty" gorm:"many2many:formal_groups;"`
	// TODO: boolean for public guestlist
	// Use pointers for updates?
}
