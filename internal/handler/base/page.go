package base

// BasePageRequest is the standard body for all paginated POST /query endpoints.
type BasePageRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortName string `json:"sort_name"`
	SortBy   string `json:"sort_by"` // "asc" | "desc"
	Search   string `json:"search"`
	IsAll    bool   `json:"is_all"` // bypass pagination when true
}

// Normalize applies defaults for missing pagination fields.
func (r *BasePageRequest) Normalize() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 || r.PageSize > 100 {
		r.PageSize = 20
	}
	if r.SortBy != "asc" && r.SortBy != "desc" {
		r.SortBy = "desc"
	}
}
