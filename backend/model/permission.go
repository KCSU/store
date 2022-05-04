package model

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `json:"-"`
	Resource  string    `json:"resource"`
	Action    string    `json:"action"`
	RoleID    uuid.UUID `json:"-"`
	Role      *Role     `json:"-"`
}
