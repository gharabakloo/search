package entity

import "strconv"

const (
	KeyPageNumber = "page_number"
	KeyPageSize   = "page_size"
)

type Pagination struct {
	CurrentPage uint64
	PerPage     uint64
}

func (p Pagination) GetOffset() uint64 {
	return (p.CurrentPage - 1) * p.PerPage
}

func (p Pagination) GetFrom() uint64 {
	return p.GetOffset() + 1
}

func (p Pagination) GetTo() uint64 {
	return p.CurrentPage * p.PerPage
}

type Page struct {
	Number string
	Size   string
}

func (p Page) Parse() Pagination {
	const (
		base            = 10
		bitSize         = 64
		defaultPageSize = 10
	)

	var pagination Pagination
	pagination.CurrentPage, _ = strconv.ParseUint(p.Number, base, bitSize)
	if pagination.CurrentPage < 1 {
		pagination.CurrentPage = 1
	}

	pagination.PerPage, _ = strconv.ParseUint(p.Size, base, bitSize)
	if pagination.PerPage < 1 {
		pagination.PerPage = defaultPageSize
	}
	return pagination
}
