package dto

type AdminGroupDto struct {
	Name   string `json:"name" validate:"required,min=3"`
	Type   string `json:"type" validate:"required,oneof=inst group manual"`
	Lookup string `json:"lookup" validate:"required_unless=Type manual"`
}
