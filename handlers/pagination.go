package handlers

import (
	"net/http"
	"strconv"
)

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
)

type PaginationMeta struct {
	Total    int  `json:"total"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	HasMore  bool `json:"has_more"`
}

func parsePagination(r *http.Request) (int, int, bool) {
	q := r.URL.Query()
	page := defaultPage
	limit := defaultLimit

	if p := q.Get("page"); p != "" {
		val, err := strconv.Atoi(p)
		if err != nil || val < 1 {
			return 0, 0, false
		}
		page = val
	}

	if l := q.Get("limit"); l != "" {
		val, err := strconv.Atoi(l)
		if err != nil || val < 1 {
			return 0, 0, false
		}
		if val > maxLimit {
			val = maxLimit
		}
		limit = val
	}

	return page, limit, true
}
