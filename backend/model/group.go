package model

import (
	"time"
)

type Group struct {
	Model
	Name       string `json:"name"`
	Type       string `json:"type"`
	Lookup     string `json:"lookup"`
	GroupUsers []GroupUser
}

// A "join table" for the User <-> Group relation.
// Note that we can't use foreign keys because the users may
// not exist.
// TODO: still maybe figure out a way to sort that out
type GroupUser struct {
	GroupID   int    `gorm:"primaryKey"`
	UserEmail string `gorm:"primaryKey"`
	CreatedAt time.Time
	Group     Group
}
