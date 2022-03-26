package dto

type UserRoleDto struct {
	UserEmail string `json:"userEmail"`
	UserName  string `json:"userName"`
	RoleName  string `json:"roleName"`
	UserID    uint   `json:"userId"`
	RoleID    uint   `json:"roleId"`
}

type AddUserRoleDto struct {
	Email  string `json:"email" validate:"required,email"`
	RoleID int    `json:"roleId"`
}

type RemoveUserRoleDto struct {
	Email  string `query:"email" validate:"required,email"`
	RoleID int    `query:"roleId"`
}
