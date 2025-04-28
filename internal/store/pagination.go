package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"min=1,max=20"`
	Offset int    `schema:"offset" validate:"min=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	// Parse limit with default value
	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	} else {
		fq.Limit = 20 // Default value
	}

	// Parse offset with default value
	if offset := qs.Get("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	} else {
		fq.Offset = 0 // Default value
	}

	// Parse sort with default value
	if sort := qs.Get("sort"); sort != "" {
		fq.Sort = sort
	} else {
		fq.Sort = "desc" // Default value
	}

	return fq, nil
}
