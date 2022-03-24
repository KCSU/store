package model

import "time"

type Permission struct {
	ID        uint      `gorm:"type:SERIAL;primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Resource  string    `json:"resource"`
	Action    string    `json:"action"`
	RoleID    uint      `json:"roleId"`
}
