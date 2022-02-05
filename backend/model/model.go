package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"type:SERIAL;primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
