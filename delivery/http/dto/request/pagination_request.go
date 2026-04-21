package request

import "github.com/Josey34/goshop/domain/valueobject"

type PaginationRequest struct {
	Page  int `form:"page" validate:"required,min=1"`
	Limit int `form:"limit" validate:"required,min=1,max=100"`
}

func (r PaginationRequest) ToPagination() valueobject.Pagination {
	page := r.Page
	if page < 1 {
		page = 1
	}

	limit := r.Limit
	if limit < 1 {
		limit = 10
	}

	return valueobject.Pagination{
		Page:  page,
		Limit: limit,
	}
}
