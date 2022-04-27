package dto

import "github.com/kcsu/store/model"

type BillStatsDto struct {
	Formals []model.FormalCostBreakdown `json:"formals"`
	Users   []model.UserCostBreakdown   `json:"users"`
}
