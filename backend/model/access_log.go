package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AccessLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	Email     string         `json:"email"`
	Message   string         `json:"message"`
	Metadata  datatypes.JSON `json:"metadata"`
}
