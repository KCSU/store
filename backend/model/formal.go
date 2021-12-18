package model

import "time"

type Formal struct {
	Model
	Name         string    `json:"name"`
	Menu         string    `json:"menu"`
	Price        float32   `json:"price" gorm:"type:decimal(10,2)"`
	GuestPrice   float32   `json:"guestPrice" gorm:"type:decimal(10,2)"`
	GuestLimit   uint      `json:"guestLimit"`
	Tickets      uint      `json:"tickets"`
	GuestTickets uint      `json:"guestTickets"`
	SaleStart    time.Time `json:"saleStart"`
	SaleEnd      time.Time `json:"saleEnd"`
	DateTime     time.Time `json:"dateTime"`
}
