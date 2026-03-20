package utils

import (
	"strconv"
)

type Pagination struct {
	Page         int32 `json:"page"`
	Limit        int32 `json:"limit"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int32 `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

func NewPagination(page int32, limit int32, totalRecords int64) *Pagination {

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		envLimit := GetEnv("LIMIT_ITEMS_ON_PER_PAGE", "10")
		limitInt, err := strconv.Atoi(envLimit)
		if err != nil || limitInt <= 0 {
			limitInt = 10
		}
		limit = int32(limitInt)
	}

	totalPages := int32((totalRecords + int64(limit) - 1) / int64(limit))
	return &Pagination{
		Page:         page,
		Limit:        limit,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		HasNext:      totalPages > page,
		HasPrev:      page > 1,
	}
}

func NewPaginationResponse(data any ,page int32, limit int32, totalRecords int64) map[string]any {
	return map[string]any {
		"data": data,
		"pagination": NewPagination(page, limit, totalRecords),
	}
}