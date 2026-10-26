package models

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

type Pagination struct {
	Page  int
	Limit int
}

func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = DefaultLimit
	}
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

type Metadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Paginated is the struct returned by every list endpoint.
type Paginated[T any] struct {
	Data     []T      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

// NewPaginated wraps a page of items with its metadata.
func NewPaginated[T any](data []T, p Pagination, total int) Paginated[T] {
	if data == nil {
		data = []T{}
	}
	pages := 0
	if p.Limit > 0 {
		pages = (total + p.Limit - 1) / p.Limit
	}
	return Paginated[T]{
		Data:     data,
		Metadata: Metadata{Page: p.Page, Limit: p.Limit, Total: total, TotalPages: pages},
	}
}
