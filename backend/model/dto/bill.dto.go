package dto

type BillDto struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Start string `json:"start" validate:"required,datetime=2006-01-02"`
	End   string `json:"end" validate:"required,datetime=2006-01-02"`
}
