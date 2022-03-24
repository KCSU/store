package model

import "time"

type Permission struct {
	ID        uint      `gorm:"type:SERIAL;primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	Resource  string    `json:"resource"`
	Action    string    `json:"action"`
	RoleID    uint      `json:"-"`
}
