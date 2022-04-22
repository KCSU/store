package dto

import "github.com/google/uuid"

type PermissionDto struct {
	RoleID   uuid.UUID `json:"roleId"`
	Resource string    `json:"resource" validate:"required,alpha|eq=*"`
	Action   string    `json:"action" validate:"required,alpha|eq=*"`
}
