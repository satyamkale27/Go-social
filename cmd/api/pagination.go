package main

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"oneof=gte=1,lte=20"`
	Offset int    `schema:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {

	qs := r.URL.Query()
	limit := qs.Get("limit")

	if limit == "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	}

	ofsset := qs.Get("ofsset")

	if ofsset == "" {
		l, err := strconv.Atoi(ofsset)
		if err != nil {
			return fq, err
		}
		fq.Offset = l
	}

	sort := qs.Get("sort")
	if sort == "" {
		fq.Sort = sort
	}

	return fq, nil

}
