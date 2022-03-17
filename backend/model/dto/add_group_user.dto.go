package dto

type AddGroupUserDto struct {
	Email string `json:"userEmail" validate:"required,email"`
}
