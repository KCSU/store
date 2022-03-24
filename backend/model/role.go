package model

type Role struct {
	Model
	Name        string
	Permissions []Permission
}
