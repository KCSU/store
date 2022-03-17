package dto

type GroupUserDto struct {
	Email string `json:"userEmail" query:"email" validate:"required,email"`
}
