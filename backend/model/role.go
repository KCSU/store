package model

type Role struct {
	Model
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}
