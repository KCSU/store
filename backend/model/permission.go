package model

import "time"

type Permission struct {
	ID        uint      `gorm:"type:SERIAL;primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Resource  string
	Action    string
	RoleID    uint
}
