package model

type Ticket struct {
	Model
	IsGuest    bool    `json:"isGuest"`
	IsQueue    bool    `json:"isQueue"`
	MealOption string  `json:"option"`
	FormalID   int     `json:"formalId"`
	Formal     *Formal `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	UserID     int     `json:"userId"`
	User       *User   `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}
