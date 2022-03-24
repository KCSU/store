package dto

import "github.com/kcsu/store/model"

type UserDto struct {
	model.User
	Groups      []GroupDto         `json:"groups"`
	Permissions []model.Permission `json:"permissions"`
}
