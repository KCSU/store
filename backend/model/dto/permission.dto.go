package dto

type PermissionDto struct {
	RoleID   int    `json:"roleId"`
	Resource string `json:"resource" validate:"required,alpha|eq=*"`
	Action   string `json:"action" validate:"required,alpha|eq=*"`
}
