package model

type ManualTicket struct {
	Model
	MealOption string  `json:"option"`
	FormalID   int     `json:"formalId"`
	Formal     *Formal `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	// TODO: should this be an enum/custom type?
	// One of: "complimentary", "ents", "standard", "guest"
	Type          string `json:"type"`
	Name          string `json:"name"`
	Justification string `json:"justification"`
	BilledTo      string `json:"billedTo"`
}
