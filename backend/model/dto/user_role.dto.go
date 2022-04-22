package dto

import "github.com/google/uuid"

type UserRoleDto struct {
	UserEmail string    `json:"userEmail"`
	UserName  string    `json:"userName"`
	RoleName  string    `json:"roleName"`
	UserID    uuid.UUID `json:"userId"`
	RoleID    uuid.UUID `json:"roleId"`
}

type AddUserRoleDto struct {
	Email  string    `json:"email" validate:"required,email"`
	RoleID uuid.UUID `json:"roleId"`
}

type RemoveUserRoleDto struct {
	Email  string    `query:"email" validate:"required,email"`
	RoleID uuid.UUID `query:"roleId"`
}
