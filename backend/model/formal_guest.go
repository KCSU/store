package model

type FormalGuest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Guests int    `json:"guests"`
}
