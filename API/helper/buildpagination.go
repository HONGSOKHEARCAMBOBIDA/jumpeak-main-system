package helper

import (
	"mysql/model"
	"mysql/request"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func NormalizePagination(pf *request.Pagination) {
	if pf.Page < 1 {
		pf.Page = 1
	}
	if pf.PageSize <= 0 {
		pf.PageSize = DefaultPageSize
	}
	if pf.PageSize > MaxPageSize {
		pf.PageSize = MaxPageSize
	}
}

func BuildPaginationMeta(pf request.Pagination, totalCount int64) *model.PaginationMetadata {
	totalPages := 0
	if pf.PageSize > 0 {
		totalPages = (int(totalCount) + pf.PageSize - 1) / pf.PageSize
	}
	return &model.PaginationMetadata{
		CurrentPage: pf.Page,
		PageSize:    pf.PageSize,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
	}
}
