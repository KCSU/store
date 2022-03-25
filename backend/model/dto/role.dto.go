package dto

type RoleDto struct {
	Name string `json:"name" validate:"required,min=3"`
}
