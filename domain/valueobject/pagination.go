package valueobject

type Pagination struct {
	Page  int
	Limit int
}

func NewPagination(page, limit int) Pagination {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	return Pagination{Page: page, Limit: limit}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}
