package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=5"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (pfq *PaginatedFeedQuery) Parse(r *http.Request) (*PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pfq, nil
		}
		pfq.Limit = l
	}

	offset := qs.Get("offset")

	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return pfq, nil
		}
		pfq.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		pfq.Sort = sort
	}

	return pfq, nil
}
