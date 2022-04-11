package model

type ManualTicket struct {
	Model
	MealOption string
	FormalID   int
	Formal     *Formal `gorm:"constraint:OnDelete:CASCADE;"`
	// TODO: should this be an enum/custom type?
	// One of: "complimentary", "ents", "standard", "guest"
	Type          string
	Name          string
	Justification string
	BilledTo      string
}
