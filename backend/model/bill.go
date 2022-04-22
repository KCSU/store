package model

import "time"

type Bill struct {
	Model
	Name    string    `json:"name"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Formals []Formal  `json:"formals"`
}
