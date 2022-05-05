package dto

type PaginationDto struct {
	Page int `query:"page" validate:"required,min=1"`
	Size int `query:"size" validate:"required,min=1,max=100"`
}
