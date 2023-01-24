package model

type Pigeonhole struct {
	Email  string `gorm:"primaryKey"`
	Number *int
}
