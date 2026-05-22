package domain

// Paginator wraps a paginated result set with metadata.
type Paginator struct {
	Data      any   `json:"data"`
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// NewPaginator constructs a Paginator, calculating total pages from total + pageSize.
func NewPaginator(data any, page, pageSize int, total int64) Paginator {
	if pageSize == 0 {
		return Paginator{Data: data, Page: page, PageSize: pageSize, Total: total, TotalPage: 1}
	}
	totalPage := int(total) / pageSize
	if total%int64(pageSize) != 0 {
		totalPage++
	}
	return Paginator{
		Data:      data,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
	}
}
